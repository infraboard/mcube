package impl

import "github.com/prometheus/client_golang/prometheus"

func NewEventCollect() *EventCollect {
	return &EventCollect{
		errCountDesc: prometheus.NewDesc(
			"save_event_error_count",
			"事件入库失败个数统计",
			[]string{},
			prometheus.Labels{"service": "maudit"},
		),
	}
}

// 收集事件指标的采集器
type EventCollect struct {
	errCountDesc *prometheus.Desc
	// 需要自己根据实践情况来维护这个变量
	errCount int
}

func (c *EventCollect) Inc() {
	c.errCount++
}

// 指标元数据注册
func (c *EventCollect) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.errCountDesc
}

// 指标的值的采集
func (c *EventCollect) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.errCountDesc, prometheus.GaugeValue, float64(c.errCount))
}
