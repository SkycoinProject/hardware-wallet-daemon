package api

import (
	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
)

//go:generate mockery -name Gatewayer -case underscore -inpkg -testonly

// Gateway is the api gateway
type Gateway struct {
	Device *deviceWallet.Device
}

// NewGateway creates a Gateway
func NewGateway(device *deviceWallet.Device) *Gateway {
	return &Gateway{
		device,
	}
}

// Gatewayer interface for Gateway methods
type Gatewayer interface {
	deviceWallet.Devicer
}
