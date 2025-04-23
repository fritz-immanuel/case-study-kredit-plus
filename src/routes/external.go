package routes

import (
	"github.com/gin-gonic/gin"

	"case-study-kredit-plus/library/data"
	"case-study-kredit-plus/src/app/external"

	"github.com/jmoiron/sqlx"
)

func RegisterExternalRoutes(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine) {
	v1 := router.Group("/external/v1")
	{
		external.RegisterRoutes(db, dataManager, router, v1)
	}
}
