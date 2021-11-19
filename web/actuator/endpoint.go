package actuator

import (
	"gitlab.id.vin/vincart/common/golib/actuator"
	"gitlab.id.vin/vincart/common/golib/web/response"
	"net/http"
)

type Endpoint struct {
	healthService actuator.HealthService
	infoService   actuator.InfoService
}

func NewEndpoint(healthService actuator.HealthService, infoService actuator.InfoService) *Endpoint {
	return &Endpoint{
		healthService: healthService,
		infoService:   infoService,
	}
}

func (c Endpoint) HealthService() actuator.HealthService {
	return c.healthService
}

func (c Endpoint) InfoService() actuator.InfoService {
	return c.infoService
}

func (c Endpoint) Health(w http.ResponseWriter, r *http.Request) {
	health := c.healthService.Check()
	var res response.Response
	if health.Status == actuator.StatusDown {
		res = response.New(http.StatusServiceUnavailable, "Server is down", health)
	} else {
		res = response.New(http.StatusOK, "Server is up", health)
	}
	response.Write(w, res)
}

func (c Endpoint) Info(w http.ResponseWriter, r *http.Request) {
	info := c.infoService.Info()
	response.Write(w, response.Ok(info))
}
