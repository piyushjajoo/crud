# crud
crud is a CLI utility which helps in scaffolding a simple go based micro-service along with 
build scripts, api documentation, micro-service documentation and k8s deployment manifests

## Install Crud CLI

Run following command to install the cli -

```shell
go install github.com/piyushjajoo/crud@latest
```

### crud help

```
crud is a CLI utility which helps in scaffolding a simple go based micro-service along with
build scripts, api documentation, micro-service documentation and k8s deployment manifests

Usage:
  crud [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  init        init creates the scaffolding for the go based micro-service

Flags:
  -h, --help   help for crud

Use "crud [command] --help" for more information about a command.
```

## Initialize Command

Initialize command initializes the project. 
If `--swagger` flag is provided, swagger documentation will be created.
If `--chart` flag is provided, helm chart to deploy the service will be created.

```shell
crud init github.com/piyushjajoo/inventory --swagger --chart
```

### crud init help

```
Init (cobra init) command initializes the go module along with a bare-bone http-server.
Please make sure you have go installed and GOPATH set. Also make sure you have helm v3 installed as well.

By default if no flags provided it initializes following -
1. go.mod and go.sum files
2. Dockerfile to build your micro-service along with build.sh script
3. main.go with bare http-server written in gorilla mux
4. README.md with basic Summary

If you want api documentation provide --swagger flag. If you want helm chart provide --chart flag.

Usage:
  crud init <module name> [flags]

Aliases:
  init, initialize, initialise, create

Flags:
  -c, --chart         to generate helm chart
  -h, --help          help for init
  -n, --name string   module name for the go module, last part of the name will be used for directory name (e.g 'github.com/piyushjajoo/crud' is the module name and crud is the directory name)
  -s, --swagger       to generate swagger api documentation file
```
