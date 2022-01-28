module mainModule

go 1.15

require github.com/gin-gonic/gin v1.7.3

require (
	// go: domainModule@v0.0.0: malformed module path "domainModule": missing dot in first path element
	domainModule v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/kisielk/gotool v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220118154757-00ab72f36ad5 // indirect
	google.golang.org/grpc v1.43.0 // indirect
)

replace domainModule => ../domainModule
