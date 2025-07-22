package chapters

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) UnhideChapter(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы"})
		return
	}

	code, err := unhideChapter(h.DB, uint(chapterID))
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "глава успешно удалена из скрытых"})
}

func unhideChapter(db *gorm.DB, chapterID uint) (code int, err error) {
	result := db.Exec("UPDATE chapters SET hidden = FALSE WHERE id = ?", chapterID)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("глава не найдена среди скрытых")
	}

	return 0, nil
}
