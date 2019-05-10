package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gogo/protobuf/proto"
	gcli "github.com/urfave/cli"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
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

				enc := json.NewEncoder(os.Stdout)
				if err = enc.Encode(features); err != nil {
					log.Errorln(err)
					return
				}
				ff := deviceWallet.NewFirmwareFeatures(uint64(*features.FirmwareFeatures))
				if err := ff.Unmarshal(); err != nil {
					log.Errorln(err)
					return
				}
				log.Printf("\n\nFirmware features:\n%s", ff)
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
