package bus

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/infraboard/mcube/bus"
	"github.com/infraboard/mcube/bus/broker/kafka"
	"github.com/infraboard/mcube/bus/broker/nats"
	"github.com/infraboard/mcube/bus/event"
	"github.com/infraboard/mcube/logger/zap"
)

var pubCmd = &cobra.Command{
	Use:   "pub",
	Short: "发布事件",
	Long:  `发布事件`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := zap.DevelopmentSetup()
		if err != nil {
			return err
		}
		log := zap.L().Named("Bus")

		var (
			pub bus.PubManager
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
	},
}

func init() {
	BusCmd.AddCommand(pubCmd)
}
