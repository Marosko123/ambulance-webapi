package ambulance_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implAmbulancesAPI struct {
}

func NewAmbulancesApi() AmbulancesAPI {
	return &implAmbulancesAPI{}
}

func (o implAmbulancesAPI) CreateAmbulance(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implAmbulancesAPI) DeleteAmbulance(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
