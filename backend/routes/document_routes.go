package routes

import (
	"net/http"
	"time"
	"vita-track-ai/models"
	"vita-track-ai/service"

	"github.com/gin-gonic/gin"
)

// @Summary Create Document
// @Tags Document
// @Router /documents [post]
func createDocument(c *gin.Context) {
	service.CreateDocument(c)
}

// @Summary Get Document
// @Tags Document
// @Router /documents/{id} [get]
func getDocument(c *gin.Context) {
	service.GetDocument(c)
}

// @Summary Delete Document
// @Tags Document
// @Router /documents/{id} [delete]
func deleteDocument(c *gin.Context) {
	service.DeleteDocument(c)
}

// @Summary Get Calendar Documents
// @Tags Document
// @Router /documents/calendar [post]
func getCalendarDocuments(c *gin.Context) {
	service.GetCalendarDocuments(c)
}

// @Summary Update Document
// @Tags Document
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Document ID"
// @Param category formData string false "Category [Prescription or Medical Report]"
// @Param document_name formData string false "DocumentName"
// @Param document_date formData string false "DocumentDate (YYYY-MM-DD)"
// @Param Tags formData string false "tags"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /documents/update/{id} [patch]
func updateDocument(c *gin.Context) {
	var updateDocReq models.UpdateDocumentRequest
	err := c.ShouldBind(&updateDocReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to pass the values into the doc object",
			"error":   err.Error(),
		})

		return
	}

	userID := c.MustGet("user_id").(int64)
	documentId := c.Param("id")

	if updateDocReq.DocumentDate != nil {
		_, err = time.Parse("2006-01-02", *updateDocReq.DocumentDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid document_date format (use YYYY-MM-DD)",
			})
			return
		}
	}

	err = service.UpdateDocument(userID, documentId, &updateDocReq)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some problem in updating the document deatails",
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document Updated Successfully",
	})

}
