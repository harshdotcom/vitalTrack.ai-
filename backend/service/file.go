package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vita-track-ai/models"
	"vita-track-ai/repository"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var MAX_ALLOWED_SIZE_FILE int64 = 5242880

func UploadFiles(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "files are required",
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no files uploaded",
		})
		return
	}

	allowed := map[string]bool{
		".pdf": true,
		".jpg": true,
		".png": true,
	}

	var response []gin.H

	for _, file := range files {

		// normalize extension
		ext := strings.ToLower(filepath.Ext(file.Filename))

		// validate extension
		if !allowed[ext] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid file type: " + file.Filename,
			})
			return
		}

		if file.Size > MAX_ALLOWED_SIZE_FILE {
			c.JSON(413, gin.H{
				"error": "File size should not be more than 5MB " + file.Filename,
			})
			return
		}

		userID := c.MustGet("user_id").(int64)
		exceedsLimit, err := exceedStorageLimit(userID, file.Size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check storage limit",
			})
			return
		}

		if exceedsLimit {
			c.JSON(413, gin.H{
				"error": "User storage limit exceeded. Please delete some files before uploading new ones.",
			})
			return
		}

		// generate unique filename (S3 object key)
		storedName := uuid.New().String() + ext

		// upload file to S3
		if err := UploadToS3(file, storedName, os.Getenv("AWS_BUCKET_NAME")); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to upload to s3",
				"error":   err.Error(),
			})
			return
		}

		// create DB model
		fileModel := models.File{
			OriginalName: file.Filename,
			StoredName:   storedName,
			S3Key:        storedName, // ✔ correct S3 key
			FileSize:     file.Size,
			MimeType:     file.Header.Get("Content-Type"),
			UploadedBy:   userID,
		}

		// save metadata to DB
		if err := repository.CreateFile(&fileModel); err != nil {

			// optional: rollback S3 upload if DB fails
			// _ = DeleteFromS3(bucket, storedName)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to save file metadata",
			})
			return
		}

		// response object
		response = append(response, gin.H{
			"file_id":       fileModel.ID,
			"original_name": fileModel.OriginalName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"files": response,
	})
}

func GetFileDownloadURL(fileID string) (string, error) {

	file, err := repository.GetFileByID(fileID)
	if err != nil {
		return "", err
	}

	bucket := os.Getenv("AWS_BUCKET_NAME")

	url, err := GenerateSignedURL(bucket, file.S3Key)
	if err != nil {
		return "", err
	}

	return url, nil
}

func GenerateOCRText(fileId string) (string, error) {

	ctx := context.Background()
	s3Key, err := repository.GetS3Key(fileId)
	if err != nil {
		return "", err
	}

	bucket := os.Getenv("AWS_BUCKET_NAME")

	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}

	client := textract.NewFromConfig(awsCfg)

	// 3️⃣ Start async Textract job
	startOut, err := client.StartDocumentTextDetection(ctx,
		&textract.StartDocumentTextDetectionInput{
			DocumentLocation: &types.DocumentLocation{
				S3Object: &types.S3Object{
					Bucket: aws.String(bucket),
					Name:   aws.String(s3Key),
				},
			},
		})

	if err != nil {
		return "", err
	}

	// 4️⃣ Wait for job completion
	var fullText strings.Builder
	var nextToken *string

	for {
		out, err := client.GetDocumentTextDetection(ctx,
			&textract.GetDocumentTextDetectionInput{
				JobId:     aws.String(*startOut.JobId),
				NextToken: nextToken,
			})
		if err != nil {
			return "", err
		}

		switch out.JobStatus {
		case types.JobStatusInProgress:
			time.Sleep(3 * time.Second)
			continue

		case types.JobStatusFailed:
			return "", fmt.Errorf("textract job failed")

		case types.JobStatusSucceeded:
			for _, block := range out.Blocks {
				if block.BlockType == types.BlockTypeLine && block.Text != nil {
					fullText.WriteString(*block.Text)
					fullText.WriteString("\n")
				}
			}

			if out.NextToken == nil {
				return fullText.String(), nil
			}

			nextToken = out.NextToken
		}
	}

}

func DeleteFile(c *gin.Context) {

	id := c.Param("id")
	userID := c.MustGet("user_id").(int64)
	storageKey, err := repository.GetS3Key(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Unable to get S3 key from DB",
			"message": err.Error(),
		})
		return
	}

	err = DeleteFileFromS3(storageKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Unable to delete from S3 bucket",
			"message": err.Error(),
		})
		return
	}

	err = repository.DeleteFile(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "file not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file deleted",
	})
}
