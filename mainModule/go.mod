module mainModule

go 1.15

require github.com/gin-gonic/gin v1.7.3

require (
	// go: domainModule@v0.0.0: malformed module path "domainModule": missing dot in first path element
	domainModule v0.0.0-00010101000000-000000000000
	github.com/go-kit/kit v0.12.0
	github.com/go-kit/log v0.2.0
	github.com/gogo/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/json-iterator/go v1.1.12
	github.com/julienschmidt/httprouter v1.3.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/segmentio/ksuid v1.0.4
	github.com/spf13/viper v1.10.1
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

replace domainModule => ../domainModule
