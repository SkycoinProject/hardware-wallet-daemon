package devicewallet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/gogo/protobuf/proto"
	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/usb"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"
)

type DeviceDriver interface {
	SendToDevice(dev io.ReadWriteCloser, chunks [][64]byte) (wire.Message, error)
	SendToDeviceNoAnswer(dev io.ReadWriteCloser, chunks [][64]byte) error
}

const (
	// DeviceTypeEmulator use emulator
	DeviceTypeEmulator = iota + 1
	// DeviceTypeUsb use usb
	DeviceTypeUSB
	// DeviceTypeInvalid not valid value
	DeviceTypeInvalid
)

type Driver struct {
	DeviceType
}

func (drv *Driver) SendToDeviceNoAnswer(dev io.ReadWriteCloser, chunks [][64]byte) error {
	for _, element := range chunks {
		_, err := dev.Write(element[:])
		if err != nil {
			return err
		}
	}
	return nil
}

func (drv *Driver) SendToDevice(dev io.ReadWriteCloser, chunks [][64]byte) (wire.Message, error) {
	var msg wire.Message
	for _, element := range chunks {
		_, err := dev.Write(element[:])
		if err != nil {
			return msg, err
		}
	}
	_, err := msg.ReadFrom(dev)
	return msg, err
}

func getEmulatorDevice() (net.Conn, error) {
	return net.Dial("udp", "127.0.0.1:21324")
}

func getUsbDevice() (usb.Device, error) {
	w, err := usb.InitWebUSB()
	if err != nil {
		log.Printf("webusb: %s", err)
		return nil, err
	}
	h, err := usb.InitHIDAPI()
	if err != nil {
		log.Printf("hidapi: %s", err)
		return nil, err
	}
	b := usb.Init(w, h)

	var infos []usb.Info
	infos, err = b.Enumerate()
	if len(infos) <= 0 {
		return nil, err
	}
	tries := 0
	for tries < 3 {
		dev, err := b.Connect(infos[0].Path)
		if err != nil {
			log.Print(err.Error())
			tries++
			time.Sleep(100 * time.Millisecond)
		} else {
			return dev, err
		}
	}
	return nil, err
}

func binaryWrite(message io.Writer, data interface{}) {
	err := binary.Write(message, binary.BigEndian, data)
	if err != nil {
		log.Panic(err)
	}
}

func makeSkyWalletMessage(data []byte, msgID messages.MessageType) [][64]byte {
	message := new(bytes.Buffer)
	binaryWrite(message, []byte("##"))
	binaryWrite(message, uint16(msgID))
	binaryWrite(message, uint32(len(data)))
	binaryWrite(message, []byte("\n"))
	if len(data) > 0 {
		binaryWrite(message, data[1:])
	}

	messageLen := message.Len()
	var chunks [][64]byte
	i := 0
	for messageLen > 0 {
		var chunk [64]byte
		chunk[0] = '?'
		copy(chunk[1:], message.Bytes()[63*i:63*(i+1)])
		chunks = append(chunks, chunk)
		messageLen -= 63
		i = i + 1
	}
	return chunks
}

func getDevice(dt DeviceType) (io.ReadWriteCloser, error) {
	var dev io.ReadWriteCloser
	var err error
	switch dt {
	case DeviceTypeEmulator:
		dev, err = getEmulatorDevice()
	case DeviceTypeUSB:
		dev, err = getUsbDevice()
	}
	if dev == nil && err == nil {
		err = errors.New("No device connected")
	}
	return dev, err
}

// Initialize send an init request to the device
func initialize(d *Device) error {
	dev, err := getDevice(d.DeviceType)
	if err != nil {
		return err
	}
	defer dev.Close()
	var chunks [][64]byte

	initialize := &messages.Initialize{}
	data, err := proto.Marshal(initialize)
	if err != nil {
		return err
	}

	chunks = makeSkyWalletMessage(data, messages.MessageType_MessageType_Initialize)
	_, err = d.Driver.SendToDevice(dev, chunks)

	return err
}

func DecodeSuccessOrFailMsg(msg wire.Message) (string, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_Success) {
		return DecodeSuccessMsg(msg)
	}
	if msg.Kind == uint16(messages.MessageType_MessageType_Failure) {
		return DecodeFailMsg(msg)
	}

	return "", fmt.Errorf("calling DecodeSuccessOrFailMsg on message kind %s", messages.MessageType(msg.Kind))
}

// DecodeSuccessMsg convert byte data into string containing the success message returned by the device
func DecodeSuccessMsg(msg wire.Message) (string, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_Success) {
		success := &messages.Success{}
		err := proto.Unmarshal(msg.Data, success)
		if err != nil {
			return "", err
		}
		return success.GetMessage(), nil
	}

	return "", fmt.Errorf("calling DecodeSuccessMsg with wrong message type: %s", messages.MessageType(msg.Kind))
}

// DecodeFailMsg convert byte data into string containing the failure returned by the device
func DecodeFailMsg(msg wire.Message) (string, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_Failure) {
		failure := &messages.Failure{}
		err := proto.Unmarshal(msg.Data, failure)
		if err != nil {
			return "", err
		}
		return failure.GetMessage(), nil
	}
	return "", fmt.Errorf("calling DecodeFailMsg with wrong message type: %s", messages.MessageType(msg.Kind))
}

// DecodeResponseSkycoinAddress convert byte data into list of addresses, meant to be used after DevicePinMatrixAck
func DecodeResponseSkycoinAddress(msg wire.Message) ([]string, error) {
	log.Printf("%x\n", msg.Data)

	if msg.Kind == uint16(messages.MessageType_MessageType_ResponseSkycoinAddress) {
		responseSkycoinAddress := &messages.ResponseSkycoinAddress{}
		err := proto.Unmarshal(msg.Data, responseSkycoinAddress)
		if err != nil {
			return []string{}, err
		}
		return responseSkycoinAddress.GetAddresses(), nil
	}

	return []string{}, fmt.Errorf("calling DecodeResponseSkycoinAddress with wrong message type: %s", messages.MessageType(msg.Kind))
}

// DecodeResponseTransactionSign convert byte data into list of signatures
func DecodeResponseTransactionSign(msg wire.Message) ([]string, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_ResponseTransactionSign) {
		responseSkycoinTransactionSign := &messages.ResponseTransactionSign{}
		err := proto.Unmarshal(msg.Data, responseSkycoinTransactionSign)
		if err != nil {
			return make([]string, 0), err
		}
		return responseSkycoinTransactionSign.GetSignatures(), nil
	}

	return []string{}, fmt.Errorf("calling DecodeResponseeSkycoinSignMessage with wrong message type: %s", messages.MessageType(msg.Kind))
}

// DecodeResponseSkycoinSignMessage convert byte data into signed message, meant to be used after DevicePinMatrixAck
func DecodeResponseSkycoinSignMessage(msg wire.Message) (string, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_ResponseSkycoinSignMessage) {
		responseSkycoinSignMessage := &messages.ResponseSkycoinSignMessage{}
		err := proto.Unmarshal(msg.Data, responseSkycoinSignMessage)
		if err != nil {
			return "", err
		}
		return responseSkycoinSignMessage.GetSignedMessage(), nil
	}
	return "", fmt.Errorf("calling DecodeResponseeSkycoinSignMessage with wrong message type: %s", messages.MessageType(msg.Kind))
}
