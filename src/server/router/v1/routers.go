package v1

type RouteInterface interface {
}

type RouteServer struct {
}

const (
	ApiVersion = "/apis/v1"
)

func NewRouteServer() RouteInterface {
	return &RouteServer{}
}

const (
	StatusOK = "ok"
)

type Object struct {
	any
}
