package dto

import "github.com/Araks1255/mangacage/pkg/common/models/dto"

type TagOnModerationDTO struct {
	dto.ResponseTagDTO

	CreatorID uint   `json:"creatorId,omitempty"`
	Creator   string `json:"creator,omitempty"`
}
