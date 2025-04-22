package user

import (
	"case-study-kredit-plus/library/types"
	"case-study-kredit-plus/models"

	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(*gin.Context, models.FindAllUserParams) ([]*models.User, *types.Error)
	Find(*gin.Context, string) (*models.User, *types.Error)
	Count(*gin.Context, models.FindAllUserParams) (int, *types.Error)
	Create(*gin.Context, models.User) (*models.User, *types.Error)
	Update(*gin.Context, string, models.User) (*models.User, *types.Error)

	FindStatus(*gin.Context) ([]*models.Status, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.User, *types.Error)

	// LOGIN
	Login(*gin.Context, models.FindAllUserParams) (*models.UserJWTContent, *types.Error)
}
