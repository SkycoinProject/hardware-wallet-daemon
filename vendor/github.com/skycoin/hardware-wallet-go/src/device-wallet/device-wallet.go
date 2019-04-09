package devicewallet

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/skycoin/skycoin/src/util/logging"

	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"
)

var (
	log = logging.MustGetLogger("device-wallet")
)

const (
	entropyBufferSize int = 32
)

// ButtonType is emulator button press simulation type
type ButtonType int32

const (
	// ButtonLeft press left button
	ButtonLeft ButtonType = iota

	// ButtonRight press right button
	ButtonRight
	// ButtonBoth press both buttons
	ButtonBoth
)

//go:generate mockery -name Devicer -case underscore -inpkg -testonly

// Devicer provides api for the hw wallet functions
type Devicer interface {
	AddressGen(addressN, startIndex int, confirmAddress bool) (wire.Message, error)
	ApplySettings(usePassphrase bool, label string) (wire.Message, error)
	Backup() (wire.Message, error)
	Cancel() (wire.Message, error)
	CheckMessageSignature(message, signature, address string) (wire.Message, error)
	ChangePin() (wire.Message, error)
	Connected() bool
	FirmwareUpload(payload []byte, hash [32]byte) error
	GetFeatures() (wire.Message, error)
	GenerateMnemonic(wordCount uint32, usePassphrase bool) (wire.Message, error)
	Recovery(wordCount uint32, usePassphrase, dryRun bool) (wire.Message, error)
	SetMnemonic(mnemonic string) (wire.Message, error)
	TransactionSign(inputs []*messages.SkycoinTransactionInput, outputs []*messages.SkycoinTransactionOutput) (wire.Message, error)
	SignMessage(addressIndex int, message string) (wire.Message, error)
	Wipe() (wire.Message, error)
	PinMatrixAck(p string) (wire.Message, error)
	WordAck(word string) (wire.Message, error)
	PassphraseAck(passphrase string) (wire.Message, error)
	ButtonAck() (wire.Message, error)
	SetAutoPressButton(simulateButtonPress bool, simulateButtonType ButtonType) error
}

// Device provides hardware wallet functions
type Device struct {
	Driver DeviceDriver

	// dev latest device connection instance
	// during an ongoing operation the device instance cannot be requested before closing the previous instance
	// keeping the connection instance in the struct helps with closing and opening of the connection
	dev io.ReadWriteCloser

	simulateButtonPress bool
	simulateButtonType  ButtonType
}

// DeviceTypeFromString returns device type from string
func DeviceTypeFromString(deviceType string) DeviceType {
	var dtRet DeviceType
	switch deviceType {
	case DeviceTypeUSB.String():
		dtRet = DeviceTypeUSB
	case DeviceTypeEmulator.String():
		dtRet = DeviceTypeEmulator
	default:
		log.Errorf("device type not set, valid options are %s or %s",
			DeviceTypeUSB,
			DeviceTypeEmulator)
		dtRet = DeviceTypeInvalid
	}
	return dtRet
}

// NewDevice returns a new device instance
func NewDevice(deviceType DeviceType) (device *Device) {
	switch deviceType {
	case DeviceTypeUSB, DeviceTypeEmulator:
		device = &Device{
			&Driver{deviceType},
			nil,
			false,
			ButtonType(-1),
		}
	default:
		device = nil
	}
	return device
}

// Connect makes a connection to the connected device
func (d *Device) Connect() error {
	// close any existing connections
	if d.dev != nil {
		d.dev.Close()
	}

	dev, err := d.Driver.GetDevice()
	if err != nil {
		return err
	}

	d.dev = dev
	return nil
}

// AddressGen Ask the device to generate an address
func (d *Device) AddressGen(addressN, startIndex int, confirmAddress bool) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()

	chunks, err := MessageAddressGen(addressN, startIndex, confirmAddress)
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, chunks)
}

// ApplySettings send ApplySettings request to the device
func (d *Device) ApplySettings(usePassphrase bool, label string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()

	chunks, err := MessageApplySettings(usePassphrase, label)
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, chunks)
}

// Backup ask the device to perform the seed backup
func (d *Device) Backup() (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()
	var msg wire.Message

	var chunks [][64]byte
	err := Initialize(d.dev)
	if err != nil {
		return wire.Message{}, err
	}

	chunks, err = MessageBackup()
	if err != nil {
		return wire.Message{}, err
	}

	msg, err = d.Driver.SendToDevice(d.dev, chunks)
	if err != nil {
		return wire.Message{}, err
	}

	for msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
		msg, err = d.ButtonAck()
		if err != nil {
			return wire.Message{}, err
		}
	}

	return msg, nil
}

// Cancel sends a Cancel request
func (d *Device) Cancel() (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()
	chunks, err := MessageCancel()
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, chunks)
}

// CheckMessageSignature Check a message signature matches the given address.
func (d *Device) CheckMessageSignature(message, signature, address string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()

	// Send CheckMessageSignature
	chunks, err := MessageCheckMessageSignature(message, signature, address)
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, chunks)
}

// ChangePin changes device's PIN code
// The message that is sent contains an encoded form of the PIN.
// The digits of the PIN are displayed in a 3x3 matrix on the Trezor,
// and the message that is sent back is a string containing the positions
// of the digits on that matrix. Below is the mapping between positions
// and characters to be sent:
// 7 8 9
// 4 5 6
// 1 2 3
// For example, if the numbers are laid out in this way on the Trezor,
// 3 1 5
// 7 8 4
// 9 6 2
// To set the PIN "12345", the positions are:
// top, bottom-right, top-left, right, top-right
// so you must send "83769".
func (d *Device) ChangePin() (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()
	chunks, err := MessageChangePin()
	if err != nil {
		return wire.Message{}, err
	}

	msg, err := d.Driver.SendToDevice(d.dev, chunks)
	if err != nil {
		return wire.Message{}, err
	}

	// Acknowledge that a button has been pressed
	if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
		msg, err = d.ButtonAck()
		if err != nil {
			return msg, err
		}
	}

	return msg, nil
}

// Connected check if a device is connected
func (d *Device) Connected() bool {
	dev, err := d.Driver.GetDevice()
	if dev == nil {
		return false
	}
	defer dev.Close()
	if err != nil {
		return false
	}

	chunks, err := MessageConnected()
	if err != nil {
		log.Error(err)
		return false
	}
	for _, element := range chunks {
		_, err = dev.Write(element[:])
		if err != nil {
			return false
		}
	}
	var msg wire.Message
	_, err = msg.ReadFrom(dev)
	if err != nil {
		return false
	}
	return msg.Kind == uint16(messages.MessageType_MessageType_Success)
}

// FirmwareUpload Updates device's firmware
func (d *Device) FirmwareUpload(payload []byte, hash [32]byte) error {
	if d.Driver.DeviceType() != DeviceTypeUSB {
		return errors.New("wrong device type")
	}
	if err := d.Connect(); err != nil {
		return err
	}
	defer d.dev.Close()

	if err := Initialize(d.dev); err != nil {
		return err
	}

	log.Printf("Length of firmware %d", uint32(len(payload)))

	chunks, err := MessageFirmwareErase(payload)
	if err != nil {
		return err
	}
	erasemsg, err := d.Driver.SendToDevice(d.dev, chunks)
	if err != nil {
		return err
	}

	switch erasemsg.Kind {
	case uint16(messages.MessageType_MessageType_Success):
		log.Printf("Success %d! FirmwareErase %s\n", erasemsg.Kind, erasemsg.Data)
	case uint16(messages.MessageType_MessageType_Failure):
		msg, err := DecodeFailMsg(erasemsg)
		if err != nil {
			return err
		}

		return errors.New(msg)
	default:
		return fmt.Errorf("received unexpected message type: %s", messages.MessageType(erasemsg.Kind))
	}

	log.Printf("Hash: %x\n", hash)

	chunks, err = MessageFirmwareUpload(payload, hash)
	if err != nil {
		return err
	}
	uploadmsg, err := d.Driver.SendToDevice(d.dev, chunks)
	if err != nil {
		return err
	}

	switch uploadmsg.Kind {
	case uint16(messages.MessageType_MessageType_Success):
		log.Printf("Success %d! FirmwareUpload %s\n", uploadmsg.Kind, uploadmsg.Data)
	case uint16(messages.MessageType_MessageType_Failure):
		msg, err := DecodeFailMsg(erasemsg)
		if err != nil {
			return err
		}

		return errors.New(msg)
	default:
		return fmt.Errorf("received unexpected message type: %s", messages.MessageType(erasemsg.Kind))
	}

	// Send ButtonAck
	chunks, err = MessageButtonAck()
	if err != nil {
		return err
	}
	return d.Driver.SendToDeviceNoAnswer(d.dev, chunks)
}

// GetFeatures send Features message to the device
func (d *Device) GetFeatures() (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()
	chunks, err := MessageGetFeatures()
	if err != nil {
		return wire.Message{}, err
	}

	return d.Driver.SendToDevice(d.dev, chunks)
}

// GenerateMnemonic Ask the device to generate a mnemonic and configure itself with it.
func (d *Device) GenerateMnemonic(wordCount uint32, usePassphrase bool) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()
	generateMnemonicChunks, err := MessageGenerateMnemonic(wordCount, usePassphrase)
	if err != nil {
		return wire.Message{}, err
	}
	msg, err := d.Driver.SendToDevice(d.dev, generateMnemonicChunks)
	if err != nil {
		return msg, err
	}

	switch msg.Kind {
	case uint16(messages.MessageType_MessageType_ButtonRequest):
		return d.ButtonAck()
	case uint16(messages.MessageType_MessageType_EntropyRequest):
		chunks, err := MessageEntropyAck(entropyBufferSize)
		if err != nil {
			return wire.Message{}, err
		}
		msg, err = d.Driver.SendToDevice(d.dev, chunks)
		if err != nil {
			return wire.Message{}, err
		}
		msg, err = d.Driver.SendToDevice(d.dev, generateMnemonicChunks)
		if err != nil {
			return msg, err
		}
	}

	return msg, err
}

// Recovery ask the device to perform the seed backup
func (d *Device) Recovery(wordCount uint32, usePassphrase, dryRun bool) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()
	var msg wire.Message
	var chunks [][64]byte

	log.Printf("Using passphrase %t\n", usePassphrase)
	chunks, err := MessageRecovery(wordCount, usePassphrase, dryRun)
	if err != nil {
		return wire.Message{}, err
	}
	msg, err = d.Driver.SendToDevice(d.dev, chunks)
	if err != nil {
		return msg, err
	}
	log.Printf("Recovery device response kind is: %d\n", msg.Kind)

	if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
		msg, err = d.ButtonAck()
		if err != nil {
			return wire.Message{}, err
		}
	}

	return msg, nil
}

// SetMnemonic Configure the device with a mnemonic.
func (d *Device) SetMnemonic(mnemonic string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()

	// Send SetMnemonic
	chunks, err := MessageSetMnemonic(mnemonic)
	if err != nil {
		return wire.Message{}, err
	}
	msg, err := d.Driver.SendToDevice(d.dev, chunks)
	if err != nil {
		return wire.Message{}, err
	}

	if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
		msg, err = d.ButtonAck()
		if err != nil {
			return wire.Message{}, err
		}
	}

	return msg, err
}

// SignMessage Ask the device to sign a message using the secret key at given index.
func (d *Device) SignMessage(addressIndex int, message string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()

	chunks, err := MessageSignMessage(addressIndex, message)
	if err != nil {
		return wire.Message{}, err
	}
	return d.Driver.SendToDevice(d.dev, chunks)
}

// TransactionSign Ask the device to sign a transaction using the given information.
func (d *Device) TransactionSign(inputs []*messages.SkycoinTransactionInput, outputs []*messages.SkycoinTransactionOutput) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()
	chunks, err := MessageTransactionSign(inputs, outputs)
	if err != nil {
		return wire.Message{}, err
	}
	return d.Driver.SendToDevice(d.dev, chunks)
}

// Wipe wipes out device configuration
func (d *Device) Wipe() (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()
	var chunks [][64]byte

	err := Initialize(d.dev)
	if err != nil {
		return wire.Message{}, err
	}

	chunks, err = MessageWipe()
	if err != nil {
		return wire.Message{}, err
	}

	var msg wire.Message
	msg, err = d.Driver.SendToDevice(d.dev, chunks)
	if err != nil {
		return wire.Message{}, err
	}
	log.Printf("Wipe device %d! Answer is: %x\n", msg.Kind, msg.Data)

	if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
		msg, err = d.ButtonAck()
		if err != nil {
			return wire.Message{}, err
		}
	}

	return msg, err
}

// ButtonAck when the device is waiting for the user to press a button
// the PC need to acknowledge, showing it knows we are waiting for a user action
func (d *Device) ButtonAck() (wire.Message, error) {
	var msg wire.Message
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()

	// Send ButtonAck
	chunks, err := MessageButtonAck()
	if err != nil {
		return msg, err
	}
	err = sendToDeviceNoAnswer(d.dev, chunks)
	if err != nil {
		return msg, err
	}

	// simulate button press
	if d.simulateButtonPress {
		if err := d.SimulateButtonPress(); err != nil {
			return msg, err
		}
	}

	_, err = msg.ReadFrom(d.dev)
	return msg, err
}

// PassphraseAck send this message when the device is waiting for the user to input a passphrase
func (d *Device) PassphraseAck(passphrase string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()
	chunks, err := MessagePassphraseAck(passphrase)
	if err != nil {
		return wire.Message{}, err
	}
	return d.Driver.SendToDevice(d.dev, chunks)
}

// WordAck send a word to the device during device "recovery procedure"
func (d *Device) WordAck(word string) (wire.Message, error) {
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()
	chunks, err := MessageWordAck(word)
	if err != nil {
		return wire.Message{}, err
	}
	msg, err := d.Driver.SendToDevice(d.dev, chunks)
	if err != nil {
		return wire.Message{}, err
	}

	return msg, nil
}

// PinMatrixAck during PIN code setting use this message to send user input to device
func (d *Device) PinMatrixAck(p string) (wire.Message, error) {
	time.Sleep(1 * time.Second)
	if err := d.Connect(); err != nil {
		return wire.Message{}, err
	}
	defer d.dev.Close()

	log.Printf("Setting pin: %s\n", p)

	chunks, err := MessagePinMatrixAck(p)
	if err != nil {
		return wire.Message{}, nil
	}
	return d.Driver.SendToDevice(d.dev, chunks)
}

// SimulateButtonPress simulates a button press on emulator
func (d *Device) SimulateButtonPress() error {
	if d.Driver.DeviceType() != DeviceTypeEmulator {
		return fmt.Errorf("wrong device type: %s", d.Driver.DeviceType())
	}

	simulateMsg, err := MessageSimulateButtonPress(d.simulateButtonType)
	if err != nil {
		return err
	}

	_, err = d.dev.Write(simulateMsg.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// SetAutoPressButton enables and sets button press type
func (d *Device) SetAutoPressButton(simulateButtonPress bool, simulateButtonType ButtonType) error {
	if d.Driver.DeviceType() == DeviceTypeEmulator {
		d.simulateButtonPress = simulateButtonPress

		if simulateButtonPress {
			switch simulateButtonType {
			case ButtonLeft, ButtonRight, ButtonBoth:
				d.simulateButtonType = simulateButtonType
			default:
				return fmt.Errorf("invalid button type: %d", simulateButtonType)
			}
		} else {
			// set invalid button press type
			d.simulateButtonType = 3
		}
	}

	return nil
}
