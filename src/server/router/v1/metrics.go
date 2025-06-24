package v1

import (
	"github.com/clarechu/infra-pulse/src/metrics"
	"github.com/emicklei/go-restful/v3"
)

func MetricsHandler() *restful.WebService {
	ws := new(restful.WebService)
	handler := metrics.NewHandler(true, 40)
	ws.Path("/metrics")
	ws.Route(ws.GET("").
		To(handler.Metrics).
		Doc("监控信息查询").
		Operation("metrics"))
	return ws
}
