package bus

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/infraboard/mcube/bus"
	"github.com/infraboard/mcube/bus/broker/kafka"
	"github.com/infraboard/mcube/bus/broker/nats"
	"github.com/infraboard/mcube/bus/event"
	"github.com/infraboard/mcube/logger/zap"
)

var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "接收事件",
	Long:  `接收事件`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := zap.DevelopmentSetup()
		if err != nil {
			return err
		}
		log := zap.L().Named("Bus")

		var (
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

		if err := sub.Connect(); err != nil {
			return fmt.Errorf("connect to bus error, %s", err)
		}

		sub.Sub(topic, func(topic string, e *event.Event) error {
			fmt.Printf("sub event: %s\n", e)
			return nil
		})

		time.Sleep(10 * time.Minute)

		return nil
	},
}

func init() {
	BusCmd.AddCommand(subCmd)
}
