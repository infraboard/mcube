package bus

import (
	"encoding/json"

	"github.com/infraboard/mcube/bus/broker/kafka"
	"github.com/infraboard/mcube/bus/broker/nats"
	"github.com/infraboard/mcube/bus/event"
	"github.com/spf13/cobra"
)

// ProjectCmd 初始化系统
var BusCmd = &cobra.Command{
	Use:   "bus",
	Short: "事件总线调试",
	Long:  `事件总线调试`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var (
	nc = nats.NewDefaultConfig()
	kc = kafka.NewDefultConfig()
)

var (
	topic       string
	contentType string
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

func init() {
	BusCmd.PersistentFlags().StringVarP(&busType, "type", "b", "nats", "bus type, options [nats/kafka]")
	BusCmd.PersistentFlags().StringArrayVarP(&servers, "servers", "s", []string{"nats://127.0.0.1:4222"}, "bus server address")
	BusCmd.PersistentFlags().StringVarP(&username, "user", "u", "", "bus auth username")
	BusCmd.PersistentFlags().StringVarP(&password, "pass", "p", "", "bus auth password")
	BusCmd.PersistentFlags().StringVarP(&topic, "topic", "t", event.Type_OPERATE.String(), "pub/sub topic name")
	BusCmd.PersistentFlags().StringVarP(&contentType, "content-type", "c", "protobuf", "body content type, options [json/protobuf]")
}
