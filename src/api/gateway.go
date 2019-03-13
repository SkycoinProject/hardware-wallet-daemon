package api

import (
	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
)

//go:generate mockery -name Gatewayer -case underscore -inpkg -testonly

// Gateway bundles both USB and Emulator device into a single object
type Gateway struct {
	USBDevice      *deviceWallet.Device
	EmulatorDevice *deviceWallet.Device
}

// NewGateway creates a Gateway
func NewGateway(usb, emu *deviceWallet.Device) *Gateway {
	return &Gateway{
		usb,
		emu,
	}
}

// Gatewayer interface for Gateway methods
type Gatewayer interface {
	deviceWallet.Devicer
}
