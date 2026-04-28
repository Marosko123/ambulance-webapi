package ambulance_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Marosko123/ambulance-webapi/internal/db_service"
)

type implAmbulancesAPI struct {
}

func NewAmbulancesApi() AmbulancesAPI {
	return &implAmbulancesAPI{}
}

func resolveAmbulanceDb(c *gin.Context) (db_service.DbService[Ambulance], bool) {
	value, exists := c.Get("db_service")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db_service not found in context"})
		return nil, false
	}
	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db_service has unexpected type"})
		return nil, false
	}
	return db, true
}

// CreateAmbulance - Saves new ambulance definition
func (o *implAmbulancesAPI) CreateAmbulance(c *gin.Context) {
	db, ok := resolveAmbulanceDb(c)
	if !ok {
		return
	}

	ambulance := Ambulance{}
	if err := c.BindJSON(&ambulance); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ambulance.Id == "" {
		ambulance.Id = uuid.New().String()
	}

	err := db.CreateDocument(c, ambulance.Id, &ambulance)

	switch err {
	case nil:
		c.JSON(http.StatusCreated, ambulance)
	case db_service.ErrConflict:
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "ambulance already exists"})
	default:
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	}
}

// DeleteAmbulance - Deletes specific ambulance
func (o *implAmbulancesAPI) DeleteAmbulance(c *gin.Context) {
	db, ok := resolveAmbulanceDb(c)
	if !ok {
		return
	}

	ambulanceId := c.Param("ambulanceId")
	err := db.DeleteDocument(c, ambulanceId)

	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "ambulance not found"})
	default:
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	}
}
