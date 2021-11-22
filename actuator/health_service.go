package actuator

import "context"

type HealthService interface {
	Check(ctx context.Context) Health
}

type DefaultHealthService struct {
	checkers []HealthChecker
}

func NewDefaultHealthService(checkers []HealthChecker) HealthService {
	return &DefaultHealthService{checkers: checkers}
}

func (h DefaultHealthService) Check(ctx context.Context) Health {
	health := Health{
		Status: StatusUp,
	}
	if len(h.checkers) > 0 {
		health.Components = make(map[string]StatusDetails)
		for _, checker := range h.checkers {
			details := checker.Check(ctx)
			if details.Status == StatusDown {
				health.Status = StatusDown
			}
			health.Components[checker.Component()] = details
		}
	}
	return health
}
