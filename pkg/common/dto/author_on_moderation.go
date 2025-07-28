package dto

import "github.com/Araks1255/mangacage/pkg/common/models/dto"

type AuthorOnModerationDTO struct {
	dto.ResponseAuthorDTO

	CreatorID uint   `json:"creatorId,omitempty"`
	Creator   string `json:"creator,omitempty"`
}
