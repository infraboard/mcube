package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rs/xid"
	"github.com/spf13/cobra"

	"github.com/infraboard/mcube/bus/broker/nats"
	"github.com/infraboard/mcube/bus/event"
)

var (
	nc = nats.NewDefaultConfig()
)

var (
	topic string
	mod   string
)

func newRandomEvent() (string, error) {
	data := &event.OperateEventData{
		Session:   "mcube bus cli",
		Account:   "mcube",
		RequestId: xid.New().String(),
		IpAddress: "127.0.0.1",
		UserAgent: "mcube/cli",
		UserName:  "mcube",
	}
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// EnumCmd 枚举生成器
var BusCmd = &cobra.Command{
	Use:   "bus",
	Short: "事件总线",
	Long:  `事件总线 客户端`,
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := nats.NewBroker(nc)
		if err != nil {
			return err
		}

		if err := b.Connect(); err != nil {
			return err
		}

		switch mod {
		case "pub":
			for {
				var eventJson string
				randomE, err := newRandomEvent()
				if err != nil {
					return err
				}
				err = survey.AskOne(
					&survey.Input{
						Message: "请输入JSON格式事件:",
						Default: randomE,
					},
					&eventJson,
					survey.WithValidator(survey.Required),
				)
				if err != nil {
					return err
				}
				oe := &event.OperateEventData{}
				err = json.Unmarshal([]byte(eventJson), oe)
				if err != nil {
					return err
				}
				e, err := event.NewOperateEvent(oe)
				if err != nil {
					return err
				}
				fmt.Println(e)
				if err := b.Pub(topic, e); err != nil {
					fmt.Println(err)
				}
				fmt.Println()
			}
		case "sub":
			b.Sub(topic, func(topic string, e *event.Event) error {
				fmt.Println(e)
				return nil
			})

			time.Sleep(10 * time.Minute)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(BusCmd)
}

func init() {
	BusCmd.PersistentFlags().StringArrayVarP(&nc.Servers, "servers", "s", []string{"nats://127.0.0.1:4222"}, "bus server address")
	BusCmd.PersistentFlags().StringVarP(&topic, "topic", "t", event.Type_Operate.String(), "pub/sub topic name")
	BusCmd.PersistentFlags().StringVarP(&mod, "mod", "m", "pub", "bus run mod, options [pub/sub]")
}
