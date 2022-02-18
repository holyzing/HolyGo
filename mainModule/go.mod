module mainModule

go 1.15

require github.com/gin-gonic/gin v1.7.3

require (
	// go: domainModule@v0.0.0: malformed module path "domainModule": missing dot in first path element
	domainModule v0.0.0-00010101000000-000000000000
	github.com/spf13/viper v1.10.1
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

replace domainModule => ../domainModule
