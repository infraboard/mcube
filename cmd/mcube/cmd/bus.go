package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/infraboard/mcube/bus"
	"github.com/infraboard/mcube/bus/broker/kafka"
	"github.com/infraboard/mcube/bus/broker/nats"
	"github.com/infraboard/mcube/bus/event"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	nc = nats.NewDefaultConfig()
	kc = kafka.NewDefultConfig()
)

var (
	topic       string
	contentType string
	mod         string
	busType     string
	username    string
	password    string
	servers     []string
)

func newRandomEvent() (string, error) {
	data := &event.OperateEventData{
		Session:  "mcube bus cli",
		Account:  "mcube",
		UserName: "mcube",
	}
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// BusCmd 枚举生成器
var BusCmd = &cobra.Command{
	Use:   "bus",
	Short: "事件总线",
	Long:  `事件总线 客户端`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := zap.DevelopmentSetup()
		if err != nil {
			return err
		}
		log := zap.L().Named("Bus")

		var (
			pub bus.PubManager
			sub bus.SubManager
		)
		switch busType {
		case "nats":
			nc.Servers = servers
			nc.Username = username
			nc.Password = password
			ins, err := nats.NewBroker(nc)
			if err != nil {
				return err
			}
			ins.Debug(log)
			pub = ins
			sub = ins
		case "kafka":
			kc.Hosts = servers
			kc.Username = username
			kc.Password = password
			kp, err := kafka.NewPublisher(kc)
			if err != nil {
				return err
			}
			kp.Debug(log)

			ks, err := kafka.NewSubscriber(kc)
			if err != nil {
				return err
			}
			ks.Debug(log)
		default:
			return fmt.Errorf("unknown bus type: %s", busType)
		}
		switch mod {
		case "pub":
			if err := pub.Connect(); err != nil {
				return fmt.Errorf("connect to bus error, %s", err)
			}

			for {
				var eventJSON string
				randomE, err := newRandomEvent()
				if err != nil {
					return err
				}
				err = survey.AskOne(
					&survey.Input{
						Message: "请输入JSON格式事件:",
						Default: randomE,
					},
					&eventJSON,
					survey.WithValidator(survey.Required),
				)
				if err != nil {
					return err
				}
				oe := &event.OperateEventData{}
				err = json.Unmarshal([]byte(eventJSON), oe)
				if err != nil {
					return err
				}
				var e *event.Event
				switch contentType {
				case "json":
					e, err = event.NewJsonOperateEvent(oe)
				default:
					e, err = event.NewProtoOperateEvent(oe)
				}

				if err != nil {
					return err
				}

				// 打印事件数据
				if err := pub.Pub(topic, e); err != nil {
					log.Errorf("pub event error, %s", err)
				}
				fmt.Println()
			}
		case "sub":
			if err := sub.Connect(); err != nil {
				return fmt.Errorf("connect to bus error, %s", err)
			}

			sub.Sub(topic, func(topic string, e *event.Event) error {
				fmt.Printf("sub event: %s\n", e)
				return nil
			})

			time.Sleep(10 * time.Minute)
		default:
			return fmt.Errorf("unknown mod: %s", mod)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(BusCmd)
}

func init() {
	BusCmd.PersistentFlags().StringVarP(&busType, "type", "b", "nats", "bus type, options [nats/kafka]")
	BusCmd.PersistentFlags().StringArrayVarP(&servers, "servers", "s", []string{"nats://127.0.0.1:4222"}, "bus server address")
	BusCmd.PersistentFlags().StringVarP(&username, "user", "u", "", "bus auth username")
	BusCmd.PersistentFlags().StringVarP(&password, "pass", "p", "", "bus auth password")
	BusCmd.PersistentFlags().StringVarP(&topic, "topic", "t", event.Type_Operate.String(), "pub/sub topic name")
	BusCmd.PersistentFlags().StringVarP(&mod, "mod", "m", "pub", "bus run mod, options [pub/sub]")
	BusCmd.PersistentFlags().StringVarP(&contentType, "content-type", "c", "protobuf", "body content type, options [json/protobuf]")
}
