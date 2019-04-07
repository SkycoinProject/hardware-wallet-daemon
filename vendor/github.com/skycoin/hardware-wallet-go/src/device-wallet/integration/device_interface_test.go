package integration

import (
	"fmt"
	"log"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"
)

func testHelperGetDeviceWithBestEffort(testName string, t *testing.T) *deviceWallet.Device {
	emDevice := deviceWallet.NewDevice(deviceWallet.DeviceTypeEmulator)
	usbDevice := deviceWallet.NewDevice(deviceWallet.DeviceTypeUSB)
	if usbDevice.Connected() {
		return usbDevice
	} else if emDevice.Connected() {
		return emDevice
	}
	t.Skip(testName + " does not work if neither Emulator nor USB device is connected")
	return nil
}

func TestDevice(t *testing.T) {
	device := testHelperGetDeviceWithBestEffort("TestDevice", t)
	if device == nil {
		return
	}
	// var msg wire.Message
	// var chunks [][64]byte
	// var inputWord string
	// var err error

	if device.Driver.DeviceType() == deviceWallet.DeviceTypeEmulator {
		err := device.SetAutoPressButton(true, deviceWallet.ButtonRight)
		require.NoError(t, err)
	}

	_, err := device.Wipe()
	require.NoError(t, err)

	// Send ChangePin message
	// chunks = MessageRecoveryDevice(12)
	// msg = SendToDevice(dev, chunks)
	// if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
	//     chunks = MessageButtonAck()
	//     msg = SendToDevice(dev, chunks)
	// }
	// for msg.Kind == uint16(messages.MessageType_MessageType_WordRequest) {
	//     fmt.Print("Word request: ")
	//     fmt.Scanln(&inputWord)
	//     chunks = MessageWordAck(strings.TrimSpace(inputWord))
	//     msg = SendToDevice(dev, chunks)
	// }
	// fmt.Printf("Response: %s\n", messages.MessageType_name[int32(msg.Kind)])
	// if msg.Kind == uint16(messages.MessageType_MessageType_Failure) {
	//     failMsg := &messages.Failure{}
	//     proto.Unmarshal(msg.Data, failMsg)
	//     fmt.Printf("Code: %d\nMessage: %s\n", failMsg.GetCode(), failMsg.GetMessage());
	// }

	_, err = device.SetMnemonic("cloud flower upset remain green metal below cup stem infant art thank")
	require.NoError(t, err)

	msg, err := device.AddressGen(9, 15, false)
	require.NoError(t, err)

	addresses, err := deviceWallet.DecodeResponseSkycoinAddress(msg)
	require.NoError(t, err)
	i := 0
	require.Equal(t, 9, len(addresses))
	require.Equal(t, addresses[i], "3NpgZ6g1UWZc5f5B7gC3hU6NhyEWxznohG")
	i++
	require.Equal(t, addresses[i], "Wr6wE5bHwBpg4kTs3EF4xi2cLs2dEWy1BN")
	i++
	require.Equal(t, addresses[i], "2DpKC15mSBhNMptvLgudRim6ScY4df1TwLd")
	i++
	require.Equal(t, addresses[i], "ZdaQWbWers3qYpKKSoBNq237CXQhGmHwX9")
	i++
	require.Equal(t, addresses[i], "9mTMfX1v6TnCYCK8frzSKAL4m2Lx1uu7Kq")
	i++
	require.Equal(t, addresses[i], "2cKu9tZz3eGqo6ny7D447o4RpMFNEk8KyXr")
	i++
	require.Equal(t, addresses[i], "2mqM8j7Zqq5MiWLEgJyAzTAPQ9sd575nh9X")
	i++
	require.Equal(t, addresses[i], "29pYKsirWo21ZPhEsdNmcCVExgAeK5ShpMF")
	i++
	require.Equal(t, addresses[i], "n6ou5D4hSGCXsAiVCJX6y6jc454xvcoSet")
	// chunks = MessageBackupDevice()
	// msg = SendToDevice(dev, chunks)
	// for msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
	//     chunks = MessageButtonAck()
	//     msg = SendToDevice(dev, chunks)
	// }
	// fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

	msg, err = device.AddressGen(1, 1, false)
	require.NoError(t, err)
	addresses, err = deviceWallet.DecodeResponseSkycoinAddress(msg)
	require.NoError(t, err)
	require.Equal(t, len(addresses), 1)
	require.Equal(t, addresses[0], "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")

	message := "Hello World!"
	msg, err = device.SignMessage(1, message)
	require.NoError(t, err)
	signature, err := deviceWallet.DecodeResponseSkycoinSignMessage(msg)
	require.NoError(t, err)
	log.Print(signature)
	require.Equal(t, 130, len(signature))

	msg, err = device.CheckMessageSignature(message, signature, addresses[0])
	require.NoError(t, err)
	require.Equal(t, "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs", string(msg.Data[2:]))
}

func TestGetAddressUsb(t *testing.T) {
	device := deviceWallet.NewDevice(deviceWallet.DeviceTypeUSB)
	if !device.Connected() {
		t.Skip("TestGetAddressUsb do not work if Usb device is not connected")
		return
	}

	_, err := device.Wipe()
	require.NoError(t, err)
	// need to connect the usb device
	_, err = device.SetMnemonic("cloud flower upset remain green metal below cup stem infant art thank")
	require.NoError(t, err)
	msg, err := device.AddressGen(2, 0, false)
	require.NoError(t, err)
	addresses, err := deviceWallet.DecodeResponseSkycoinAddress(msg)
	require.NoError(t, err)
	log.Print(addresses)
	require.Equal(t, addresses[0], "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.Equal(t, addresses[1], "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")
}

func TestGetAddressEmulator(t *testing.T) {
	device := deviceWallet.NewDevice(deviceWallet.DeviceTypeEmulator)
	if !device.Connected() {
		t.Skip("TestGetAddressEmulator do not work if emulator is not running")
		return
	}

	err := device.SetAutoPressButton(true, deviceWallet.ButtonRight)
	require.NoError(t, err)

	_, err = device.Wipe()
	require.NoError(t, err)

	_, err = device.SetMnemonic("cloud flower upset remain green metal below cup stem infant art thank")
	require.NoError(t, err)

	msg, err := device.AddressGen(2, 0, false)
	require.NoError(t, err)
	addresses, err := deviceWallet.DecodeResponseSkycoinAddress(msg)
	require.NoError(t, err)
	log.Print(addresses)
	require.Equal(t, addresses[0], "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.Equal(t, addresses[1], "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")
}

func TransactionToDevice(deviceType deviceWallet.DeviceType, transactionInputs []*messages.SkycoinTransactionInput, transactionOutputs []*messages.SkycoinTransactionOutput) (wire.Message, error) {
	device := deviceWallet.NewDevice(deviceType)
	if device == nil {
		return wire.Message{}, fmt.Errorf("invalid device type: %s", deviceType)
	}

	if device.Driver.DeviceType() == deviceWallet.DeviceTypeEmulator {
		err := device.SetAutoPressButton(true, deviceWallet.ButtonRight)
		if err != nil {
			return wire.Message{}, err
		}
	}

	msg, err := device.TransactionSign(transactionInputs, transactionOutputs)
	if err != nil {
		return wire.Message{}, err
	}
	for {
		switch msg.Kind {
		case uint16(messages.MessageType_MessageType_ResponseTransactionSign):
			return msg, nil
		case uint16(messages.MessageType_MessageType_Success):
			fmt.Println("Should end with ResponseTransactionSign request")
		case uint16(messages.MessageType_MessageType_ButtonRequest):
			msg, err = device.ButtonAck()
			if err != nil {
				return wire.Message{}, err
			}
		case uint16(messages.MessageType_MessageType_PassphraseRequest):
			var passphrase string
			fmt.Printf("Input passphrase: ")
			fmt.Scanln(&passphrase)
			msg, err = device.PassphraseAck(passphrase)
			if err != nil {
				return wire.Message{}, err
			}
		case uint16(messages.MessageType_MessageType_PinMatrixRequest):
			var pinEnc string
			fmt.Printf("PinMatrixRequest response: ")
			fmt.Scanln(&pinEnc)
			msg, err = device.PinMatrixAck(pinEnc)
			if err != nil {
				return wire.Message{}, err
			}
		case uint16(messages.MessageType_MessageType_Failure):
			failMsg, err := deviceWallet.DecodeFailMsg(msg)
			if err != nil {
				return wire.Message{}, err
			}
			fmt.Printf("Failed with message: %s\n", failMsg)
		default:
			return wire.Message{}, fmt.Errorf("received unexpected message type: %s", messages.MessageType(msg.Kind))
		}
	}
}

func TestTransactions(t *testing.T) {
	device := testHelperGetDeviceWithBestEffort("TestTransactions", t)
	if device == nil {
		return
	}

	if device.Driver.DeviceType() == deviceWallet.DeviceTypeEmulator {
		err := device.SetAutoPressButton(true, deviceWallet.ButtonRight)
		require.NoError(t, err)
	}

	_, err := device.Wipe()
	require.NoError(t, err)

	_, err = device.SetMnemonic("cloud flower upset remain green metal below cup stem infant art thank")
	require.NoError(t, err)

	var transactionInputs []*messages.SkycoinTransactionInput
	var transactionOutputs []*messages.SkycoinTransactionOutput
	var transactionInput messages.SkycoinTransactionInput
	var transactionOutput messages.SkycoinTransactionOutput
	var transactionInput1 messages.SkycoinTransactionInput
	var transactionOutput1 messages.SkycoinTransactionOutput
	var transactionInput2 messages.SkycoinTransactionInput
	var transactionOutput2 messages.SkycoinTransactionOutput

	// Sample 1
	transactionInput.HashIn = proto.String("181bd5656115172fe81451fae4fb56498a97744d89702e73da75ba91ed5200f9")
	transactionInput.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput)

	transactionOutput.Address = proto.String("K9TzLrgqz7uXn3QJHGxmzdRByAzH33J2ot")
	transactionOutput.Coin = proto.Uint64(100000)
	transactionOutput.Hour = proto.Uint64(2)
	transactionOutputs = append(transactionOutputs, &transactionOutput)

	msg, err := TransactionToDevice(device.Driver.DeviceType(), transactionInputs, transactionOutputs)
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_ResponseTransactionSign), msg.Kind)

	signatures, err := deviceWallet.DecodeResponseTransactionSign(msg)
	require.NoError(t, err)
	require.Equal(t, len(signatures), 1)

	msg, err = device.CheckMessageSignature(
		"d11c62b1e0e9abf629b1f5f4699cef9fbc504b45ceedf0047ead686979498218", signatures[0],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	transactionOutputs = transactionOutputs[:0]
	transactionInputs = transactionInputs[:0]

	// Sample 2
	// --------------------
	transactionInput.HashIn = proto.String("01a9ef6c25271229ef9760e1536c3dc5ccf0ead7de93a64c12a01340670d87e9")
	transactionInput.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput)
	transactionInput1.HashIn = proto.String("8c2c97bfd34e0f0f9833b789ce03c2e80ac0b94b9d0b99cee6ea76fb662e8e1c")
	transactionInput1.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput1)

	transactionOutput.Address = proto.String("K9TzLrgqz7uXn3QJHGxmzdRByAzH33J2ot")
	transactionOutput.Coin = proto.Uint64(20800000)
	transactionOutput.Hour = proto.Uint64(255)
	transactionOutputs = append(transactionOutputs, &transactionOutput)

	msg, err = TransactionToDevice(device.Driver.DeviceType(), transactionInputs, transactionOutputs)
	require.NoError(t, err)
	require.Equal(t, msg.Kind, uint16(messages.MessageType_MessageType_ResponseTransactionSign))

	signatures, err = deviceWallet.DecodeResponseTransactionSign(msg)
	require.NoError(t, err)
	require.Equal(t, len(signatures), len(transactionInputs))

	msg, err = device.CheckMessageSignature(
		"9bbde062d665a8b11ae15aee6d4f32f0f3d61af55160c142060795a219378a54", signatures[0],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	msg, err = device.CheckMessageSignature(
		"f947b0352b19672f7b7d04dc2f1fdc47bc5355878f3c47a43d4d4cfbae07d026", signatures[1],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	transactionOutputs = transactionOutputs[:0]
	transactionInputs = transactionInputs[:0]

	// Sample 3
	// --------------------
	transactionInput.HashIn = proto.String("da3b5e29250289ad78dc42dcf007ab8f61126198e71e8306ff8c11696a0c40f7")
	transactionInput.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput)
	transactionInput1.HashIn = proto.String("33e826d62489932905dd936d3edbb74f37211d68d4657689ed4b8027edcad0fb")
	transactionInput1.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput1)
	transactionInput2.HashIn = proto.String("668f4c144ad2a4458eaef89a38f10e5307b4f0e8fce2ade96fb2cc2409fa6592")
	transactionInput2.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput2)

	transactionOutput.Address = proto.String("K9TzLrgqz7uXn3QJHGxmzdRByAzH33J2ot")
	transactionOutput.Coin = proto.Uint64(111000000)
	transactionOutput.Hour = proto.Uint64(6464556)
	transactionOutputs = append(transactionOutputs, &transactionOutput)
	transactionOutput1.Address = proto.String("2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8")
	transactionOutput1.Coin = proto.Uint64(1900000)
	transactionOutput1.Hour = proto.Uint64(1)
	transactionOutputs = append(transactionOutputs, &transactionOutput1)

	msg, err = TransactionToDevice(device.Driver.DeviceType(), transactionInputs, transactionOutputs)
	require.NoError(t, err)
	require.Equal(t, msg.Kind, uint16(messages.MessageType_MessageType_ResponseTransactionSign))

	signatures, err = deviceWallet.DecodeResponseTransactionSign(msg)
	require.NoError(t, err)
	require.Equal(t, len(signatures), len(transactionInputs))
	msg, err = device.CheckMessageSignature(
		"ff383c647551a3ba0387f8334b3f397e45f9fc7b3b5c3b18ab9f2b9737bce039", signatures[0],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	msg, err = device.CheckMessageSignature(
		"c918d83d8d3b1ee85c1d2af6885a0067bacc636d2ebb77655150f86e80bf4417", signatures[1],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	msg, err = device.CheckMessageSignature(
		"0e827c5d16bab0c3451850cc6deeaa332cbcb88322deea4ea939424b072e9b97", signatures[2],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	transactionOutputs = transactionOutputs[:0]
	transactionInputs = transactionInputs[:0]

	// Sample 4
	// --------------------
	transactionInput.HashIn = proto.String("b99f62c5b42aec6be97f2ca74bb1a846be9248e8e19771943c501e0b48a43d82")
	transactionInput.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput)
	transactionInput1.HashIn = proto.String("cd13f705d9c1ce4ac602e4c4347e986deab8e742eae8996b34c429874799ebb2")
	transactionInput1.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput1)

	transactionOutput.Address = proto.String("22S8njPeKUNJBijQjNCzaasXVyf22rWv7gF")
	transactionOutput.Coin = proto.Uint64(23100000)
	transactionOutput.Hour = proto.Uint64(0)
	transactionOutputs = append(transactionOutputs, &transactionOutput)

	msg, err = TransactionToDevice(device.Driver.DeviceType(), transactionInputs, transactionOutputs)
	require.NoError(t, err)
	require.Equal(t, msg.Kind, uint16(messages.MessageType_MessageType_ResponseTransactionSign))

	signatures, err = deviceWallet.DecodeResponseTransactionSign(msg)
	require.NoError(t, err)
	require.Equal(t, len(signatures), len(transactionInputs))
	msg, err = device.CheckMessageSignature(
		"42a26380399172f2024067a17704fceda607283a0f17cb0024ab7a96fc6e4ac6", signatures[0],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	msg, err = device.CheckMessageSignature(
		"5e0a5a8c7ea4a2a500c24e3a4bfd83ef9f74f3c2ff4bdc01240b66a41e34ebbf", signatures[1],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	transactionOutputs = transactionOutputs[:0]
	transactionInputs = transactionInputs[:0]

	// Sample 5
	// --------------------
	transactionInput.HashIn = proto.String("4c12fdd28bd580989892b0518f51de3add96b5efb0f54f0cd6115054c682e1f1")
	transactionInput.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput)

	transactionOutput.Address = proto.String("2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8")
	transactionOutput.Coin = proto.Uint64(1000000)
	transactionOutput.Hour = proto.Uint64(0)
	transactionOutputs = append(transactionOutputs, &transactionOutput)

	msg, err = TransactionToDevice(device.Driver.DeviceType(), transactionInputs, transactionOutputs)
	require.NoError(t, err)
	require.Equal(t, msg.Kind, uint16(messages.MessageType_MessageType_ResponseTransactionSign))

	signatures, err = deviceWallet.DecodeResponseTransactionSign(msg)
	require.NoError(t, err)
	require.Equal(t, len(signatures), len(transactionInputs))
	msg, err = device.CheckMessageSignature(
		"c40e110f5e460532bfb03a5a0e50262d92d8913a89c87869adb5a443463dea69", signatures[0],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	transactionOutputs = transactionOutputs[:0]
	transactionInputs = transactionInputs[:0]

	// Sample 6
	// --------------------
	transactionInput.HashIn = proto.String("c5467f398fc3b9d7255d417d9ca208c0a1dfa0ee573974a5fdeb654e1735fc59")
	transactionInput.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput)

	transactionOutput.Address = proto.String("K9TzLrgqz7uXn3QJHGxmzdRByAzH33J2ot")
	transactionOutput.Coin = proto.Uint64(10000000)
	transactionOutput.Hour = proto.Uint64(1)
	transactionOutputs = append(transactionOutputs, &transactionOutput)
	transactionOutput1.Address = proto.String("VNz8LR9JTSoz5o7qPHm3QHj4EiJB6LV18L")
	transactionOutput1.Coin = proto.Uint64(5500000)
	transactionOutput1.Hour = proto.Uint64(0)
	transactionOutputs = append(transactionOutputs, &transactionOutput1)
	transactionOutput2.Address = proto.String("22S8njPeKUNJBijQjNCzaasXVyf22rWv7gF")
	transactionOutput2.Coin = proto.Uint64(4500000)
	transactionOutput2.Hour = proto.Uint64(1)
	transactionOutputs = append(transactionOutputs, &transactionOutput2)

	msg, err = TransactionToDevice(device.Driver.DeviceType(), transactionInputs, transactionOutputs)
	require.NoError(t, err)
	require.Equal(t, msg.Kind, uint16(messages.MessageType_MessageType_ResponseTransactionSign))

	signatures, err = deviceWallet.DecodeResponseTransactionSign(msg)
	require.NoError(t, err)
	require.Equal(t, len(signatures), len(transactionInputs))

	msg, err = device.CheckMessageSignature(
		"7edea77354eca0999b1b023014eb04638b05313d40711707dd03a9935696ccd1", signatures[0],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	transactionOutputs = transactionOutputs[:0]
	transactionInputs = transactionInputs[:0]

	// Sample 7
	// --------------------
	transactionInput.HashIn = proto.String("7b65023cf64a56052cdea25ce4fa88943c8bc96d1ab34ad64e2a8b4c5055087e")
	transactionInput.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput)
	transactionInput1.HashIn = proto.String("0c0696698cba98047bc042739e14839c09bbb8bb5719b735bff88636360238ad")
	transactionInput1.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput1)
	transactionInput2.HashIn = proto.String("ae3e0b476b61734e590b934acb635d4ad26647bc05867cb01abd1d24f7f2ce50")
	transactionInput2.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput2)

	transactionOutput.Address = proto.String("22S8njPeKUNJBijQjNCzaasXVyf22rWv7gF")
	transactionOutput.Coin = proto.Uint64(25000000)
	transactionOutput.Hour = proto.Uint64(33)
	transactionOutputs = append(transactionOutputs, &transactionOutput)

	msg, err = TransactionToDevice(device.Driver.DeviceType(), transactionInputs, transactionOutputs)
	require.NoError(t, err)
	require.Equal(t, msg.Kind, uint16(messages.MessageType_MessageType_ResponseTransactionSign))

	signatures, err = deviceWallet.DecodeResponseTransactionSign(msg)
	require.NoError(t, err)
	require.Equal(t, len(signatures), len(transactionInputs))

	msg, err = device.CheckMessageSignature(
		"ec9053ab9988feb0cfb3fcce96f02c7d146ff7a164865c4434d1dbef42a24e91", signatures[0],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	msg, err = device.CheckMessageSignature(
		"332534f92c27b31f5b73d8d0c7dde4527b540024f8daa965fe9140e97f3c2b06", signatures[1],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	msg, err = device.CheckMessageSignature(
		"63f955205ceb159415268bad68acaae6ac8be0a9f33ef998a84d1c09a8b52798", signatures[2],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	transactionOutputs = transactionOutputs[:0]
	transactionInputs = transactionInputs[:0]

	// Sample 8
	// --------------------
	transactionInput.HashIn = proto.String("ae6fcae589898d6003362aaf39c56852f65369d55bf0f2f672bcc268c15a32da")
	transactionInput.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput)

	transactionOutput.Address = proto.String("3pXt9MSQJkwgPXLNePLQkjKq8tsRnFZGQA")
	transactionOutput.Coin = proto.Uint64(1000000)
	transactionOutput.Hour = proto.Uint64(1000)
	transactionOutputs = append(transactionOutputs, &transactionOutput)

	msg, err = TransactionToDevice(device.Driver.DeviceType(), transactionInputs, transactionOutputs)
	require.NoError(t, err)
	require.Equal(t, msg.Kind, uint16(messages.MessageType_MessageType_ResponseTransactionSign))

	signatures, err = deviceWallet.DecodeResponseTransactionSign(msg)
	require.NoError(t, err)
	require.Equal(t, len(signatures), len(transactionInputs))

	msg, err = device.CheckMessageSignature(
		"47bfa37c79f7960df8e8a421250922c5165167f4c91ecca5682c1106f9010a7f", signatures[0],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	transactionOutputs = transactionOutputs[:0]
	transactionInputs = transactionInputs[:0]

	// Sample 9
	// --------------------
	transactionInput.HashIn = proto.String("ae6fcae589898d6003362aaf39c56852f65369d55bf0f2f672bcc268c15a32da")
	transactionInput.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput)

	transactionOutput.Address = proto.String("3pXt9MSQJkwgPXLNePLQkjKq8tsRnFZGQA")
	transactionOutput.Coin = proto.Uint64(300000)
	transactionOutput.Hour = proto.Uint64(500)
	transactionOutputs = append(transactionOutputs, &transactionOutput)
	transactionOutput1.Address = proto.String("S6Dnv6gRTgsHCmZQxjN7cX5aRjJvDvqwp9")
	transactionOutput1.Coin = proto.Uint64(700000)
	transactionOutput1.Hour = proto.Uint64(500)
	transactionOutputs = append(transactionOutputs, &transactionOutput1)

	msg, err = TransactionToDevice(device.Driver.DeviceType(), transactionInputs, transactionOutputs)
	require.NoError(t, err)
	require.Equal(t, msg.Kind, uint16(messages.MessageType_MessageType_ResponseTransactionSign))

	signatures, err = deviceWallet.DecodeResponseTransactionSign(msg)
	require.NoError(t, err)
	require.Equal(t, len(signatures), len(transactionInputs))

	msg, err = device.CheckMessageSignature(
		"e0c6e4982b1b8c33c5be55ac115b69be68f209c5d9054954653e14874664b57d", signatures[0],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))

	transactionOutputs = transactionOutputs[:0]
	transactionInputs = transactionInputs[:0]

	// Sample 10
	// --------------------
	transactionInput.HashIn = proto.String("ae6fcae589898d6003362aaf39c56852f65369d55bf0f2f672bcc268c15a32da")
	transactionInput.Index = proto.Uint32(0)
	transactionInputs = append(transactionInputs, &transactionInput)

	transactionOutput.Address = proto.String("S6Dnv6gRTgsHCmZQxjN7cX5aRjJvDvqwp9")
	transactionOutput.Coin = proto.Uint64(1000000)
	transactionOutput.Hour = proto.Uint64(1000)
	transactionOutputs = append(transactionOutputs, &transactionOutput)

	msg, err = TransactionToDevice(device.Driver.DeviceType(), transactionInputs, transactionOutputs)
	require.NoError(t, err)
	require.Equal(t, msg.Kind, uint16(messages.MessageType_MessageType_ResponseTransactionSign))

	signatures, err = deviceWallet.DecodeResponseTransactionSign(msg)
	require.NoError(t, err)
	require.Equal(t, len(signatures), len(transactionInputs))

	msg, err = device.CheckMessageSignature(
		"457648543755580ad40ab461bbef2b0ffe19f2130f2f220cbb2f196b05d436b4", signatures[0],
		"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.NoError(t, err)
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), msg.Kind) // Success message
	require.Equal(t, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", string(msg.Data[2:]))
}
