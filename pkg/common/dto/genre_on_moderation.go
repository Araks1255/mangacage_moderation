package dto

import "github.com/Araks1255/mangacage/pkg/common/models/dto"

type GenreOnModerationDTO struct {
	dto.ResponseGenreDTO

	CreatorID uint   `json:"creatorId"`
	Creator   string `json:"creator"`
}
