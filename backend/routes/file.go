package routes

import (
	"net/http"
	"vita-track-ai/service"

	"github.com/gin-gonic/gin"
)

// @Summary Upload File
// @Router /api/v1/files/upload [post]
func uploadFile(c *gin.Context) {
	service.UploadFiles(c)
}

// @Summary Get File
// @Router /api/v1/files/{id} [get]
func getFile(c *gin.Context) {

	id := c.Param("id")

	url, err := service.GetFileDownloadURL(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate download url",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

// @Summary Get File Text
// @Router /api/v1/files/ocr/{id} [get]
func getFileText(c *gin.Context) {
	fileId := c.Param("id")
	text, err := service.GenerateOCRText(fileId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "failed to generate text",
			"errorMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"text": text,
	})
}

// @Summary Get File AI Analysis
// @Router /api/v1/files/ai/{id} [get]
func getFileAnalysis(c *gin.Context) {
	fileId := c.Param("id")
	jsonText, err := service.AnalyzeMedicalReport(fileId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "failed to generate ai analysis",
			"errorMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"json": jsonText,
	})

}

// @Summary Delete File
// @Router /api/v1/files/{id} [delete]
func deleteFile(c *gin.Context) {
	service.DeleteFile(c)
}
