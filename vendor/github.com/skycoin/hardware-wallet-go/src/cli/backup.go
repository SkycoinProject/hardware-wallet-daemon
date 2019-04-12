package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
)

func backupCmd() gcli.Command {
	name := "backup"
	return gcli.Command{
		Name:         name,
		Usage:        "Ask the device to perform the seed backup procedure.",
		Description:  "",
		OnUsageError: onCommandUsageError(name),
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:   "deviceType",
				Usage:  "Device type to send instructions to, hardware wallet (USB) or emulator.",
				EnvVar: "DEVICE_TYPE",
			},
		},
		Action: func(c *gcli.Context) {
			device := deviceWallet.NewDevice(deviceWallet.DeviceTypeFromString(c.String("deviceType")))
			if device == nil {
				return
			}

			msg, err := device.Backup()
			if err != nil {
				log.Error(err)
				return
			}

			if msg.Kind == uint16(messages.MessageType_MessageType_PinMatrixRequest) {
				var pinEnc string
				fmt.Printf("PinMatrixRequest response: ")
				fmt.Scanln(&pinEnc)
				msg, err := device.PinMatrixAck(pinEnc)
				if err != nil {
					log.Error(err)
					return
				}

				for msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
					msg, err = device.ButtonAck()
					if err != nil {
						log.Error(err)
						return
					}
				}
			}

			responseMsg, err := deviceWallet.DecodeSuccessOrFailMsg(msg)
			if err != nil {
				log.Error(err)
				return
			}

			fmt.Println(responseMsg)
		},
	}
}
