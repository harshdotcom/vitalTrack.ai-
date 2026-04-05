package routes

import (
	"net/http"
	"vita-track-ai/models"
	"vita-track-ai/service"

	"github.com/gin-gonic/gin"
)

// @Summary Create Document
// @Router /api/v1/documents [post]
func createDocument(c *gin.Context) {
	service.CreateDocument(c)
}

// @Summary Get Document
// @Router /api/v1/documents/{id} [get]
func getDocument(c *gin.Context) {
	service.GetDocument(c)
}

// @Summary Delete Document
// @Router /api/v1/documents/{id} [delete]
func deleteDocument(c *gin.Context) {
	service.DeleteDocument(c)
}

// @Summary Get Calendar Documents
// @Router /api/v1/documents/calendar [post]
func getCalendarDocuments(c *gin.Context) {
	service.GetCalendarDocuments(c)
}

// @Summary Update Document
// @Router /api/v1/documents/update/{id} [patch]
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
