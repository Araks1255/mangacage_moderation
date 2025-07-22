package chapters

import (
	"errors"
	"log"
	"strconv"

	mongoModels "github.com/Araks1255/mangacage/pkg/common/models/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) GetChapterOnModerationPage(c *gin.Context) {
	chapterOnModerationID, page, err := parseGetChapterOnModerationPageParams(c.Param)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"chapter_on_moderation_id": chapterOnModerationID}
	projection := bson.M{"pages": bson.M{"$slice": []int{page, 1}}}
	opts := options.FindOne().SetProjection(projection)

	var result mongoModels.ChapterOnModerationPages

	if err := h.ChaptersPages.FindOne(c.Request.Context(), filter, opts).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.AbortWithStatusJSON(404, gin.H{"error": "глава на модерации не найдена"})
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(result.Pages) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "страница главы на модерации не найдена"})
		return
	}

	c.Data(200, "image/jpeg", result.Pages[0])
}

func parseGetChapterOnModerationPageParams(paramFn func(string) string) (chapterID uint64, page int, err error) {
	if chapterID, err = strconv.ParseUint(paramFn("id"), 10, 64); err != nil {
		return 0, 0, err
	}

	if page, err = strconv.Atoi(paramFn("page")); err != nil {
		return 0, 0, err
	}

	return chapterID, page, nil
}
