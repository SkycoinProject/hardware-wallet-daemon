package cli

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	gcli "github.com/urfave/cli"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
)

func featuresCmd() gcli.Command {
	name := "features"
	return gcli.Command{
		Name:         name,
		Usage:        "Ask the device Features.",
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

			msg, err := device.GetFeatures()
			if err != nil {
				log.Error(err)
				return
			}

			switch msg.Kind {
			case uint16(messages.MessageType_MessageType_Features):
				features := &messages.Features{}
				err = proto.Unmarshal(msg.Data, features)
				if err != nil {
					log.Error(err)
					return
				}

				fmt.Println(features)
			// TODO: figure out if this method can even return success or failure msg.
			case uint16(messages.MessageType_MessageType_Failure), uint16(messages.MessageType_MessageType_Success):
				msgData, err := deviceWallet.DecodeSuccessOrFailMsg(msg)
				if err != nil {
					log.Error(err)
					return
				}

				fmt.Println(msgData)
			default:
				log.Errorf("received unexpected message type: %s", messages.MessageType(msg.Kind))
			}
		},
	}
}
