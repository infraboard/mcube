package metric

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetric(path string) *Metric {
	return &Metric{
		Path: path,
	}
}

type Metric struct {
	Path string
}

func (m *Metric) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(m.Path)
	ws.Route(ws.GET("/").To(func(r *restful.Request, w *restful.Response) {
		// 基于标准库 包装了一层
		promhttp.Handler().ServeHTTP(w, r.Request)
	}))
	return ws
}
