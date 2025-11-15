// این فایل می‌تواند بخشی از article_handler.go باشد یا یک فایل جدید

package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	app_content "narrative-architecture/apps/backend/internal/application/content"
)

type ContentHandler struct {
	getChapter16UC *app_content.GetChapter16UseCase
}

func (h *ContentHandler) GetChapter16(c *gin.Context) {
	userID, _ := c.Get("userID") // از AuthMiddleware

	response, err := h.getChapter16UC.Execute(c.Request.Context(), userID.(string), c.ClientIP())
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	
	// برای جلوگیری از کش شدن محتوای حساس در مرورگر یا CDNها
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	
	c.JSON(http.StatusOK, response)
}