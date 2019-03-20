package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
)

func cancelCmd() gcli.Command {
	name := "cancel"
	return gcli.Command{
		Name:         name,
		Usage:        "Ask the device to cancel the ongoing procedure.",
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

			msg, err := device.Cancel()
			if err != nil {
				log.Error(err)
				return
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
