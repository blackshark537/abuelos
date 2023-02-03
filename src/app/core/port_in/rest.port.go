package portin

type ApiGateway interface {
	ForRoot(port string) error
}

type ApiPort struct{}

var apiGateway ApiGateway = nil

func InjectApi(gateway ApiGateway) {
	if apiGateway == nil {
		apiGateway = gateway
	}
}

func (apiPort *ApiPort) ForRoot(port string) error {
	return apiGateway.ForRoot(port)
}
