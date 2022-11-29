module github.com/openinfradev/tks-client

go 1.19

require (
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/matryer/is v1.4.0
	github.com/openinfradev/tks-proto v0.0.6-0.20221117013032-f3e8aa863671
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.13.0
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-openapi/errors v0.20.3 // indirect
	github.com/go-openapi/strfmt v0.21.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
	github.com/rivo/uniseg v0.4.2 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	go.mongodb.org/mongo-driver v1.10.2 // indirect
	golang.org/x/net v0.0.0-20220927171203-f486391704dc // indirect
	golang.org/x/sys v0.0.0-20220927170352-d9d178bc13c6 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220927151529-dcaddaf36704 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/openinfradev/tks-client => ./

//replace github.com/openinfradev/tks-proto => ../tks-proto
