module mainModule

go 1.15

require github.com/gin-gonic/gin v1.7.3

require (
	// go: domainModule@v0.0.0: malformed module path "domainModule": missing dot in first path element
	domainModule v0.0.0-00010101000000-000000000000
	golang.org/x/tools v0.1.5
)

replace domainModule => ../domainModule
