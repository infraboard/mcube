package service

const (
	APP_NAME = "HelloService"
)

type HelloService interface {
	Hello(request *HelloRequest, response *HelloResponse) error
}

type HelloRequest struct {
	MyName string `json:"my_name"`
}

type HelloResponse struct {
	Message string `json:"message"`
}
