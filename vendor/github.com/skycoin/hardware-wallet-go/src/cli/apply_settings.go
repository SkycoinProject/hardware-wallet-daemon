package cli

import (
	"fmt"

	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"

	gcli "github.com/urfave/cli"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
)

func applySettingsCmd() gcli.Command {
	name := "applySettings"
	return gcli.Command{
		Name:        name,
		Usage:       "Apply settings.",
		Description: "",
		Flags: []gcli.Flag{
			gcli.BoolFlag{
				Name:  "usePassphrase",
				Usage: "Configure a passphrase",
			},
			gcli.StringFlag{
				Name:  "label",
				Usage: "Configure a device label",
			},
			gcli.StringFlag{
				Name:   "deviceType",
				Usage:  "Device type to send instructions to, hardware wallet (USB) or emulator.",
				EnvVar: "DEVICE_TYPE",
			},
			gcli.StringFlag{
				Name:  "language",
				Usage: "Configure a device language",
				Value: "",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) {
			passphrase := c.Bool("usePassphrase")
			label := c.String("label")
			language := c.String("language")

			device := deviceWallet.NewDevice(deviceWallet.DeviceTypeFromString(c.String("deviceType")))
			if device == nil {
				return
			}

			var msg wire.Message
			msg, err := device.ApplySettings(passphrase, label, language)
			if err != nil {
				log.Error(err)
				return
			}

			for msg.Kind != uint16(messages.MessageType_MessageType_Failure) && msg.Kind != uint16(messages.MessageType_MessageType_Success) {
				if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
					msg, err = device.ButtonAck()
					if err != nil {
						log.Error(err)
						return
					}
					continue
				}

				if msg.Kind == uint16(messages.MessageType_MessageType_PinMatrixRequest) {
					var pinEnc string
					fmt.Printf("PinMatrixRequest response: ")
					fmt.Scanln(&pinEnc)
					pinAckResponse, err := device.PinMatrixAck(pinEnc)
					if err != nil {
						log.Error(err)
						return
					}
					log.Infof("PinMatrixAck response: %s", pinAckResponse)
					continue
				}
			}

			if msg.Kind == uint16(messages.MessageType_MessageType_Failure) {
				failMsg, err := deviceWallet.DecodeFailMsg(msg)
				if err != nil {
					log.Error(err)
					return
				}
				fmt.Println("Failed with code: ", failMsg)
				return
			}

			if msg.Kind == uint16(messages.MessageType_MessageType_Success) {
				successMsg, err := deviceWallet.DecodeSuccessMsg(msg)
				if err != nil {
					log.Error(err)
					return
				}
				fmt.Println("Success with code: ", successMsg)
				return
			}
		},
	}
}
