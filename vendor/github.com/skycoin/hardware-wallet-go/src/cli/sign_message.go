package cli

import (
	"fmt"

	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"

	gcli "github.com/urfave/cli"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
)

func signMessageCmd() gcli.Command {
	name := "signMessage"
	return gcli.Command{
		Name:        name,
		Usage:       "Ask the device to sign a message using the secret key at given index.",
		Description: "",
		Flags: []gcli.Flag{
			gcli.IntFlag{
				Name:  "addressN",
				Value: 0,
				Usage: "Index of the address that will issue the signature. Assume 0 if not set.",
			},
			gcli.StringFlag{
				Name:  "message",
				Usage: "The message that the signature claims to be signing.",
			},
			gcli.StringFlag{
				Name:   "deviceType",
				Usage:  "Device type to send instructions to, hardware wallet (USB) or emulator.",
				EnvVar: "DEVICE_TYPE",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) {
			device := deviceWallet.NewDevice(deviceWallet.DeviceTypeFromString(c.String("deviceType")))
			if device == nil {
				return
			}

			addressN := c.Int("addressN")
			message := c.String("message")
			var signature string

			var msg wire.Message
			msg, err := device.SignMessage(addressN, message)
			if err != nil {
				log.Error(err)
				return
			}
			for msg.Kind != uint16(messages.MessageType_MessageType_ResponseSkycoinSignMessage) && msg.Kind != uint16(messages.MessageType_MessageType_Failure) {
				if msg.Kind == uint16(messages.MessageType_MessageType_PinMatrixRequest) {
					var pinEnc string
					fmt.Printf("PinMatrixRequest response: ")
					fmt.Scanln(&pinEnc)
					msg, err = device.PinMatrixAck(pinEnc)
					if err != nil {
						log.Error(err)
						return
					}
					continue
				}

				if msg.Kind == uint16(messages.MessageType_MessageType_PassphraseRequest) {
					var passphrase string
					fmt.Printf("Input passphrase: ")
					fmt.Scanln(&passphrase)
					msg, err = device.PassphraseAck(passphrase)
					if err != nil {
						log.Error(err)
						return
					}
					continue
				}
			}

			if msg.Kind == uint16(messages.MessageType_MessageType_ResponseSkycoinSignMessage) {
				signature, err = deviceWallet.DecodeResponseSkycoinSignMessage(msg)
				if err != nil {
					log.Error(err)
					return
				}
				fmt.Printf("Success %d! the signature is: %s\n", msg.Kind, signature)
			} else {
				failMsg, err := deviceWallet.DecodeFailMsg(msg)
				if err != nil {
					log.Error(err)
					return
				}

				fmt.Printf("Failed with message: %s\n", failMsg)
				return
			}
		},
	}
}
