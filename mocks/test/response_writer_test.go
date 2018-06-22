package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/mocks"
)

func TestResponseWriter(t *testing.T) {
	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	context.JSON(500, gin.H{"error": "Internal server error"})
	response.Assert(500, "{\"error\":\"Internal server error\"}")
}
