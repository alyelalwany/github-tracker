package model

import "time"

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type EchoRequest struct {
	Message string `json:"message"`
}

type EchoResponse struct {
	Message    string    `json:"message"`
	ReceivedAt time.Time `json:"received_at"`
}
