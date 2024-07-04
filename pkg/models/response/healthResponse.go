package response

const (
	HealthStatusStarting = "STARTING"
	HealthStatusDown     = "DOWN"
	HealthStatusUp       = "UP"
)

type HealthValueStatus struct {
	Status string `json:"status"`
}

var AppHealthState = &HealthStateResponse{
	Status:     HealthStatusDown,
	Liveliness: HealthValueStatus{Status: HealthStatusDown},
	Readiness:  HealthValueStatus{Status: HealthStatusDown},
}

type HealthStateResponse struct {
	Status     string            `json:"status"`
	Liveliness HealthValueStatus `json:"liveliness"`
	Readiness  HealthValueStatus `json:"readiness"`
}
