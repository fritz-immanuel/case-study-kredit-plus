package consumertransaction

import (
	"case-study-kredit-plus/library/types"
	"case-study-kredit-plus/models"

	"github.com/gin-gonic/gin"
)

// Repository is the contract between Repository and usecase
type Repository interface {
	FindAll(*gin.Context, models.FindAllConsumerTransactionParams) ([]*models.ConsumerTransaction, *types.Error)
	Find(*gin.Context, string) (*models.ConsumerTransaction, *types.Error)
	Count(*gin.Context, models.FindAllConsumerTransactionParams) (int, *types.Error)
	Create(*gin.Context, *models.ConsumerTransaction) (*models.ConsumerTransaction, *types.Error)
	Update(*gin.Context, *models.ConsumerTransaction) (*models.ConsumerTransaction, *types.Error)

	FindStatus(*gin.Context) ([]*models.Status, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.ConsumerTransaction, *types.Error)
}
