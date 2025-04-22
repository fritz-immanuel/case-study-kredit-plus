package routes

import (
	"case-study-kredit-plus/src/app/businessweb"

	"github.com/gin-gonic/gin"

	"case-study-kredit-plus/library/data"

	"github.com/jmoiron/sqlx"
)

func RegisterWebRoutes(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine) {
	v1 := router.Group("/web/v1")
	{
		businessweb.RegisterRoutes(db, dataManager, router, v1)
	}
}
