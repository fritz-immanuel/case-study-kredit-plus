package businessweb

import (
	http_consumer "case-study-kredit-plus/src/app/businessweb/consumer"
	http_user "case-study-kredit-plus/src/app/businessweb/user"

	"case-study-kredit-plus/library/data"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	consumerHandler http_consumer.ConsumerHandler
	userHandler     http_user.UserHandler
)

func RegisterRoutes(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	v1 := v.Group("")
	{
		consumerHandler.RegisterAPI(db, dataManager, router, v1)
		userHandler.RegisterAPI(db, dataManager, router, v1)
	}
}
