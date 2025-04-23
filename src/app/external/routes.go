package external

import (
	http_consumertransaction "case-study-kredit-plus/src/app/external/consumertransaction"

	"case-study-kredit-plus/library/data"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	consumertransactionHandler http_consumertransaction.ConsumerTransactionHandler
)

func RegisterRoutes(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	v1 := v.Group("")
	{
		consumertransactionHandler.RegisterAPI(db, dataManager, router, v1)
	}
}
