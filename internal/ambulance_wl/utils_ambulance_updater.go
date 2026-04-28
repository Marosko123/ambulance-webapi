package ambulance_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Marosko123/ambulance-webapi/internal/db_service"
)

type ambulanceUpdater = func(
	ctx *gin.Context,
	ambulance *Ambulance,
) (updatedAmbulance *Ambulance, responseContent interface{}, status int)

func updateAmbulanceFunc(ctx *gin.Context, updater ambulanceUpdater) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db_service not found in context"})
		return
	}

	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db_service has unexpected type"})
		return
	}

	ambulanceId := ctx.Param("ambulanceId")

	ambulance, err := db.FindDocument(ctx, ambulanceId)

	switch err {
	case nil:
		// continue
	case db_service.ErrNotFound:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "ambulance not found"})
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	updatedAmbulance, responseObject, status := updater(ctx, ambulance)

	if updatedAmbulance != nil {
		err = db.UpdateDocument(ctx, ambulanceId, updatedAmbulance)
	}

	switch err {
	case nil:
		if responseObject != nil {
			ctx.JSON(status, responseObject)
		} else {
			ctx.AbortWithStatus(status)
		}
	case db_service.ErrNotFound:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "ambulance was deleted during update"})
	default:
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	}
}
