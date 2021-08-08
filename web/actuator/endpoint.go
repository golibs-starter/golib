package actuator

import (
	"gitlab.id.vin/vincart/golib/web/response"
	"net/http"
)

type Endpoint struct {
	healthService HealthService
	infoService   InfoService
}

func NewEndpoint(healthService HealthService, infoService InfoService) *Endpoint {
	return &Endpoint{
		healthService: healthService,
		infoService:   infoService,
	}
}

func (c Endpoint) HealthService() HealthService {
	return c.healthService
}

func (c Endpoint) InfoService() InfoService {
	return c.infoService
}

func (c Endpoint) Health(w http.ResponseWriter, r *http.Request) {
	health := c.healthService.Check()
	var res response.Response
	if health.Status == StatusDown {
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
