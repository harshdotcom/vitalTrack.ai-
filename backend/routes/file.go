package routes

import (
	"net/http"
	"vita-track-ai/service"

	"github.com/gin-gonic/gin"
)

func uploadFile(c *gin.Context) {
	service.UploadFiles(c)
}

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

	// respBody := []byte(jsonText)
	// report, err := utility.ParseResponse(respBody)

	c.JSON(http.StatusOK, gin.H{
		"json": jsonText,
	})

}
