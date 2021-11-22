package actuator

import "context"

const (
	StatusUp   = "UP"
	StatusDown = "DOWN"
)

// Health is a model represents for service health status.
// Includes Status UP or DOWN (see StatusUp, StatusDown),
// and a map of Components under test (with key is component name, value is StatusDetails)
type Health struct {
	Status     string                   `json:"status"`
	Components map[string]StatusDetails `json:"components,omitempty"`
}

// StatusDetails is a model represents for component's status.
// Includes status UP or DOWN (see StatusUp, StatusDown),
// and the reason if it's down
type StatusDetails struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// HealthChecker is an interface that implemented by component
type HealthChecker interface {

	// Component returns the name of
	// component under test
	Component() string

	// Check status of current component,
	// Return status details includes status: UP or DOWN,
	// and the reason if it's down
	Check(ctx context.Context) StatusDetails
}
