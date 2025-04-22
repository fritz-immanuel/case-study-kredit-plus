package consumer

import (
	"case-study-kredit-plus/library/types"
	"case-study-kredit-plus/models"

	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(*gin.Context, models.FindAllConsumerParams) ([]*models.Consumer, *types.Error)
	Find(*gin.Context, string) (*models.Consumer, *types.Error)
	Count(*gin.Context, models.FindAllConsumerParams) (int, *types.Error)
	Create(*gin.Context, models.Consumer) (*models.Consumer, *types.Error)
	Update(*gin.Context, string, models.Consumer) (*models.Consumer, *types.Error)

	FindStatus(*gin.Context) ([]*models.Status, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.Consumer, *types.Error)
}
