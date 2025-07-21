package dto

import "github.com/Araks1255/mangacage/pkg/common/models/dto"

type TitleOnModerationDTO struct {
	dto.ResponseTitleDTO

	CreatorID *uint   `json:"creatorId,omitempty"`
	Creator   *string `json:"creator,omitempty"`
}
