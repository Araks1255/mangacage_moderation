package views

import "github.com/gin-gonic/gin"

func (h handler) ShowMainPage(c *gin.Context) {
	c.HTML(200, "main_page.html", nil)
}