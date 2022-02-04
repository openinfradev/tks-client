module github.com/openinfradev/tks-client

go 1.16

require (
	github.com/golang/mock v1.6.0
	github.com/kr/text v0.2.0 // indirect
	github.com/matryer/is v1.4.0
	github.com/openinfradev/tks-contract v0.1.1-0.20210928021110-fe2b666327cc
	github.com/openinfradev/tks-proto v0.0.6-0.20211015003551-ed8f9541f40d
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
	golang.org/x/sys v0.0.0-20211023085530-d6a326fbbf70 // indirect
	google.golang.org/genproto v0.0.0-20211021150943-2b146023228c // indirect
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/openinfradev/tks-client => ./
