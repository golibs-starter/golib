package example

// ==================================================
// ===== Example for register a health component ====
// ==================================================

import (
	"context"
	"gitlab.id.vin/vincart/golib/actuator"
)

// NewSampleHealthChecker
// Use golib.ProvideHealthChecker(NewSampleHealthChecker) to register a health checker.
// In this example, the `/actuator/health` endpoint with return:
//{
//  "meta": {
//    "code": 200,
//    "message": "Server is up"
//  },
//  "data": {
//    "status": "UP",
//    "components": {
//      "sample": {
//        "status": "UP"
//      }
//    }
//  }
//}
func NewSampleHealthChecker() actuator.HealthChecker {
	return &SampleHealthChecker{}
}

type SampleHealthChecker struct {
}

func (h SampleHealthChecker) Component() string {
	return "sample"
}

func (h SampleHealthChecker) Check(ctx context.Context) actuator.StatusDetails {
	// Do something to check status and
	// return actuator.StatusUp or actuator.StatusDown
	return actuator.StatusDetails{
		Status: actuator.StatusUp,
	}
}
