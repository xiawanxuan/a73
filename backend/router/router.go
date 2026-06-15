package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ocean-survey-backend/handlers"
	"ocean-survey-backend/services"
	"ocean-survey-backend/websocket"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, hub *websocket.Hub) {
	pointCloudService := services.NewPointCloudService(db)
	pointCloudHandler := handlers.NewPointCloudHandler(pointCloudService)

	labelService := services.NewTerrainLabelService(db)
	labelHandler := handlers.NewLabelHandler(labelService, hub)

	annotationService := services.NewAnnotationService(db)
	annotationHandler := handlers.NewAnnotationHandler(annotationService)

	r.GET("/ws", websocket.WsHandler(hub))

	v1 := r.Group("/api/v1")
	{
		pointClouds := v1.Group("/point-clouds")
		{
			pointClouds.GET("", pointCloudHandler.List)
			pointClouds.POST("", pointCloudHandler.Create)
			pointClouds.GET("/:id", pointCloudHandler.GetByID)
			pointClouds.PUT("/:id", pointCloudHandler.Update)
			pointClouds.DELETE("/:id", pointCloudHandler.Delete)

			pointClouds.GET("/:pcID/annotations", annotationHandler.ListByPointCloud)
			pointClouds.POST("/:pcID/annotations", annotationHandler.Create)
		}

		labels := v1.Group("/labels")
		{
			labels.GET("", labelHandler.List)
			labels.POST("", labelHandler.Create)
			labels.GET("/:id", labelHandler.GetByID)
			labels.PUT("/:id", labelHandler.Update)
			labels.DELETE("/:id", labelHandler.Delete)
		}

		annotations := v1.Group("/annotations")
		{
			annotations.GET("/:id", annotationHandler.GetByID)
			annotations.PUT("/:id", annotationHandler.Update)
			annotations.DELETE("/:id", annotationHandler.Delete)
			annotations.POST("/:id/rollback", annotationHandler.Rollback)
			annotations.GET("/:id/snapshots", annotationHandler.ListSnapshots)
		}
	}
}
