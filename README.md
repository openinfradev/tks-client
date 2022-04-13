# TKS Client

**TKS**는 **Taco Kubernetes Service**의 약자로, SK Telecom이 만든 GitOps기반의 서비스 시스템을 의미합니다. `tks-client`는 TKS 쿠버네티스 클러스터와 서비스에 대한 관리기능을 제공하는 CLI 툴로써 TKS API서비스와 gRPC 프로토콜을 사용하여 통신합니다. 

```
$ tks
   ______ __ __ ____  _____ __ _            __
  /_  __// //_// __/ / ___// /(_)___  ___  / /_
   / /  / ,<  _\ \  / /__ / // // -_)/ _ \/ __/
  /_/  /_/|_|/___/  \___//_//_/ \__//_//_/\__/

TKS Client is CLI client for using TKS Service.
For more: https://github.com/openinfradev/tks-client/

Usage:
  tks [command]

Available Commands:
  cluster     Operation for TKS Cluster
  completion  generate the autocompletion script for the specified shell
  endpoint    Operation for Thanos Endpoint
  help        Help about any command
  service     Operation for TKS Service

Flags:
      --config string   config file (default is $HOME/.tks-client.yaml)
  -h, --help            help for tks
  -t, --toggle          Help message for toggle

Use "tks [command] --help" for more information about a command.
```

`tks-client`의 매뉴얼은 다음 [링크](https://openinfradev.github.io/releases/cli/overview/)에서 확인하실 수 있습니다.

# Installation
`tks`는 Linux, MAC, Windows 환경에 설치할 수 있습니다. tks github repo에서 바이너리를 다운로드 받거나 소스파일을 직접 빌드하여 설치합니다.

## Github Repo에서 다운로드

`tks` github repo의 [릴리즈 페이지](https://github.com/openinfradev/tks-client/releases/tag/v1.0.0-rc1)에서 사용하려는 시스템에 맞는 바이너리를 다운로드합니다.

```
$ VERSION=1.0.0-rc1
$ wget https://github.com/openinfradev/tks-client/releases/download/v${VERSION}/tks-client_${VERSION}_Linux_x86_64.tar.gz
$ tar xvzf tks-client_${VERSION}_Linux_x86_64.tar.gz
LICENSE
README.md
tks
$ ./tks
   ______ __ __ ____  _____ __ _            __
  /_  __// //_// __/ / ___// /(_)___  ___  / /_
   / /  / ,<  _\ \  / /__ / // // -_)/ _ \/ __/
  /_/  /_/|_|/___/  \___//_//_/ \__//_//_/\__/

TKS Client is CLI client for using TKS Service.
For more: https://github.com/openinfradev/tks-client/

Usage:
  tks [command]

Available Commands:
  cluster     Operation for TKS Cluster
  completion  generate the autocompletion script for the specified shell
  endpoint    Operation for Thanos Endpoint
  help        Help about any command
  service     Operation for TKS Service

Flags:
      --config string   config file (default is $HOME/.tks-client.yaml)
  -h, --help            help for tks
  -t, --toggle          Help message for toggle
  -v, --verbose         verbose output

Use "tks [command] --help" for more information about a command.

```

## 소스 빌드

`tks`는  [task](https://taskfile.dev/)를 사용하여 소스 빌드를 할 수 있습니다.
```
$ git clone git@github.com:openinfradev/tks-client.git
$ cd tks-client/
$ task build
task: [build] GOFLAGS=-mod=mod go build -o bin/tks main.go
$ ./bin/tks
   ______ __ __ ____  _____ __ _            __
  /_  __// //_// __/ / ___// /(_)___  ___  / /_
   / /  / ,<  _\ \  / /__ / // // -_)/ _ \/ __/
  /_/  /_/|_|/___/  \___//_//_/ \__//_//_/\__/

TKS Client is CLI client for using TKS Service.
For more: https://github.com/openinfradev/tks-client/

Usage:
  tks [command]

Available Commands:
  cluster     Operation for TKS Cluster
  completion  generate the autocompletion script for the specified shell
  endpoint    Operation for Thanos Endpoint
  help        Help about any command
  service     Operation for TKS Service

Flags:
      --config string   config file (default is $HOME/.tks-client.yaml)
  -h, --help            help for tks
  -t, --toggle          Help message for toggle
  -v, --verbose         verbose output

Use "tks [command] --help" for more information about a command.
```

# Configuration & Commands

```
$ tee ~/.tks-client.yaml << EOF
tksInfoUrl: "tks-info.taco.com:9110"
tksContractUrl: "tks-contract.taco.com:9110"
tksClusterLcmUrl: "tks-cluster-lcm.taco.com:9110"
EOF

$ tks cluster list
Using config file: /home/ubuntu/.tks-client.yaml
 NAME             ID                                     STATUS  CREATED_AT           UPDATED_AT
 test-20220411-2  76fc9f70-75fe-40b8-a82e-9a78402e53c4  RUNNING  2022-04-11 06:16:14  2022-04-11 06:30:28
 test-20220405-1  15436172-8d62-4d38-a37a-655b8cbeadd4  RUNNING  2022-04-05 05:03:42  2022-04-05 05:19:26
```

`tks-client`의 자세한 설정방법과 커맨드 사용법은 [매뉴얼](https://openinfradev.github.io/releases/cli/overview/)을 참고하시기 바랍니다.


