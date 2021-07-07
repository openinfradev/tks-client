module github.com/sktelecom/tks-client

go 1.16

require (
	github.com/golang/mock v1.5.0
	github.com/sktelecom/tks-contract v0.1.0
	github.com/sktelecom/tks-info v0.0.0-20210604043948-98579a23bbca
	github.com/sktelecom/tks-proto v0.0.5-0.20210601073957-185e6457787e
	google.golang.org/grpc v1.38.0
)

replace github.com/sktelecom/tks-client => ./
