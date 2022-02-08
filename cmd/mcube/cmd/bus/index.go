package bus

import (
	"encoding/json"

	"github.com/infraboard/mcube/bus/broker/kafka"
	"github.com/infraboard/mcube/bus/broker/nats"
	"github.com/infraboard/mcube/pb/event"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
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
	Cmd.PersistentFlags().StringVarP(&busType, "type", "b", "nats", "bus type, options [nats/kafka]")
	Cmd.PersistentFlags().StringArrayVarP(&servers, "servers", "s", []string{"nats://127.0.0.1:4222"}, "bus server address")
	Cmd.PersistentFlags().StringVarP(&username, "user", "u", "", "bus auth username")
	Cmd.PersistentFlags().StringVarP(&password, "pass", "p", "", "bus auth password")
	Cmd.PersistentFlags().StringVarP(&topic, "topic", "t", event.Type_OPERATE.String(), "pub/sub topic name")
	Cmd.PersistentFlags().StringVarP(&contentType, "content-type", "c", "protobuf", "body content type, options [json/protobuf]")
}
