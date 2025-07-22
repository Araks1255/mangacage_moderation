package titles

import (
	"errors"
	"log"
	"strconv"

	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) UnhideTitle(c *gin.Context) {
	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	code, err := unhideTitle(tx, uint(titleID))
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "тайтл успешно удален из скрытых"})
}

func unhideTitle(db *gorm.DB, titleID uint) (code int, err error) {
	result := db.Exec("UPDATE titles SET hidden = FALSE WHERE id = ?", titleID)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("тайтл не найден среди скрытых")
	}

	err = db.Exec("UPDATE chapters SET hidden = FALSE WHERE title_id = ?", titleID).Error

	if err != nil {
		return 500, err
	}

	return 0, nil
}
