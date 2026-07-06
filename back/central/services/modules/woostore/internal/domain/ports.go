package domain

import "context"

type PowerState struct {
	InstanceID string `json:"instance_id"`
	State      string `json:"state"`
	PublicIP   string `json:"public_ip,omitempty"`
	StoreURL   string `json:"store_url,omitempty"`
}

type IPowerManager interface {
	Start(ctx context.Context) (*PowerState, error)
	Stop(ctx context.Context) (*PowerState, error)
	Status(ctx context.Context) (*PowerState, error)
}
