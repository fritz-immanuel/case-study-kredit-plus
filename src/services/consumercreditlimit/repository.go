package consumercreditlimit

import (
	"case-study-kredit-plus/library/types"
	"case-study-kredit-plus/models"

	"github.com/gin-gonic/gin"
)

// Repository is the contract between Repository and usecase
type Repository interface {
	FindAll(*gin.Context, models.FindAllConsumerCreditLimitParams) ([]*models.ConsumerCreditLimit, *types.Error)
	Find(*gin.Context, string) (*models.ConsumerCreditLimit, *types.Error)
	Count(*gin.Context, models.FindAllConsumerCreditLimitParams) (int, *types.Error)
	Create(*gin.Context, *models.ConsumerCreditLimit) (*models.ConsumerCreditLimit, *types.Error)
	Update(*gin.Context, *models.ConsumerCreditLimit) (*models.ConsumerCreditLimit, *types.Error)

	FindStatus(*gin.Context) ([]*models.Status, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.ConsumerCreditLimit, *types.Error)

	// Check Credit Limit
	CheckCreditLimitAvailability(ctx *gin.Context, consumerID string, tenor int) (float64, *types.Error) 
}
