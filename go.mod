module github.com/sktelecom/tks-client

go 1.16

require (
	github.com/golang/mock v1.6.0
	github.com/sktelecom/tks-contract v0.1.1-0.20210604023929-73ffc015c1f1
	github.com/sktelecom/tks-info v0.0.0-20210722013555-def433889881
	github.com/sktelecom/tks-proto v0.0.6-0.20210622012523-ded9f951101f
	google.golang.org/grpc v1.38.0
)

replace github.com/sktelecom/tks-client => ./
