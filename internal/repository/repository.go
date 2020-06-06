package repository

import (
	"context"

	"github.com/yyh-gl/gocr/internal/sender"
)

type (
	Repository interface {
		Name() string
		SenderID() string
		FetchCodeReviewRequests(ctx context.Context) (CodeReviewRequests, error)
	}

	CodeReviewRequest interface {
		ConvertToMsgMaterial() (sender.Material, error)
	}

	CodeReviewRequests interface {
		ConvertToMsgMaterials() (sender.Materials, error)
	}
)
