package skwallet

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/wire"
)

type devicerSuit struct {
	suite.Suite
}

func (suite *devicerSuit) SetupTest() {
}

func TestDevicerSuitSuit(t *testing.T) {
	suite.Run(t, new(devicerSuit))
}

type testHelperCloseableBuffer struct {
	bytes.Buffer
}

func (cwr testHelperCloseableBuffer) Read(p []byte) (n int, err error) {
	return 0, nil
}
func (cwr testHelperCloseableBuffer) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (cwr testHelperCloseableBuffer) Close(disconnect bool) error {
	return nil
}

func (suite *devicerSuit) TestGenerateMnemonic() {
	// NOTE: Giving
	driverMock := &MockDeviceDriver{}
	driverMock.On("GetDevice").Return(&testHelperCloseableBuffer{}, nil)
	driverMock.On("SendToDevice", mock.Anything, mock.Anything).Return(
		wire.Message{Kind: uint16(messages.MessageType_MessageType_EntropyRequest), Data: nil}, nil)
	device := Device{driverMock, nil, false, ButtonType(-1)}

	// NOTE: When
	msg, err := device.GenerateMnemonic(12, false)

	// NOTE: Assert
	suite.Nil(err)
	driverMock.AssertCalled(suite.T(), "GetDevice")
	driverMock.AssertNumberOfCalls(suite.T(), "SendToDevice", 3)
	mock.AssertExpectationsForObjects(suite.T(), driverMock)
	spew.Dump(msg)
}

func (suite *devicerSuit) TestInterfacesImplemented() {
	var _ Devicer = (*Device)(nil)
	var _ DeviceDriver = (*Driver)(nil)
}
