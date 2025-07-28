package views

import (
	"github.com/gin-gonic/gin"
)

type handler struct{}

func RegisterRoutes(r *gin.Engine) {
	h := handler{}

	r.Static("/static", "static")
	r.LoadHTMLGlob("html/*.html")

	r.GET("/moderation", h.ShowMainPage)
	r.GET("/moderation/titles-on-moderation", h.ShowTitlesOnModerationPool)
	r.GET("/moderation/titles-on-moderation/reviewing-by-me", h.ShowTitlesOnModerationReviewingByMeCatalog)
	r.GET("/moderation/chapters-on-moderation", h.ShowChaptersOnModerationPool)
	r.GET("/moderation/chapters-on-moderation/reviewing-by-me", h.ShowChaptersOnModerationReviewingByMeCatalog)
	r.GET("/moderation/authors-on-moderation", h.ShowAuthorsOnModerationPool)
	r.GET("/moderation/authors-on-moderation/reviewing-by-me", h.ShowAuthorsOnModerationReviewingByMeCatalog)
	r.GET("/moderation/teams-on-moderation", h.ShowTeamsOnModerationPool)
	r.GET("/moderation/teams-on-moderation/reviewing-by-me", h.ShowTeamsOnModerationReviewingByMeCatalog)
	r.GET("/moderation/users-on-verification", h.ShowUsersOnVerificationPool)
	r.GET("/moderation/users/on-verification/reviewing-by-me", h.ShowUsersOnVerificationReviewingByMeCatalog)
	r.GET("/moderation/users/profile-changes", h.ShowUsersProfileChangesPool)
	r.GET("/moderation/users/profile-changes/reviewing-by-me", h.ShowUsersProfileChangesReviewingByMeCatalog)
	r.GET("/moderation/genres-on-moderation", h.ShowGenresOnModerationPool)
	r.GET("/moderation/genres-on-moderation/reviewing-by-me", h.ShowGenresOnModerationReviewingByMeCatalog)
	r.GET("/moderation/tags-on-moderation", h.ShowTagsOnModerationPool)
	r.GET("/moderation/tags-on-moderation/reviewing-by-me", h.ShowTagsOnModerationReviewingByMeCatalog)
	r.GET("/moderation/chapters-on-moderation/:id", h.ShowChapterOnModerationPage)
}
