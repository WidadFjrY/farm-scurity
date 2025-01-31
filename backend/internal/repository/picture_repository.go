package repository

import (
	"context"

	"gorm.io/gorm"
)

type PictureRepository interface {
	Save(ctx context.Context, tx *gorm.DB, filePath string)
}
