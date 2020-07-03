module github.com/mchmarny/dapr-event-subscriber-template

go 1.14

require (
	github.com/cloudevents/sdk-go/v2 v2.1.0
	github.com/dapr/go-sdk v0.8.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/stretchr/testify v1.6.1
	go.opencensus.io v0.22.4 // indirect
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	google.golang.org/genproto v0.0.0-20200702021140-07506425bd67 // indirect
	google.golang.org/grpc v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace github.com/dapr/go-sdk => github.com/mchmarny/go-sdk v0.8.1-0.20200701155538-5be80e65a8b1
