module github.com/openinfradev/tks-client

go 1.19

require (
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/google/uuid v1.3.0
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/openinfradev/tks-api v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.6.1
	golang.org/x/crypto v0.7.0
	sigs.k8s.io/yaml v1.2.0
)

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/go-openapi/errors v0.20.3 // indirect
	github.com/go-openapi/strfmt v0.21.3 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.mongodb.org/mongo-driver v1.10.2 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/term v0.6.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/openinfradev/tks-client => ./

//replace github.com/openinfradev/tks-api => ../tks-api
