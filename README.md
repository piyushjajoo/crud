# crud
crud is a CLI utility which helps in scaffolding a simple go based micro-service along with 
build scripts, api documentation, micro-service documentation and k8s deployment manifests

## Install Crud CLI

Run following command to install the cli -

```shell
go install github.com/piyushjajoo/crud
```

## Initialize Command

Initialize command initializes the project. 
If `--swagger` flag is provided, swagger documentation will be created.
If `--chart` flag is provided, helm chart to deploy the service will be created.

```shell
crud init github.com/piyushjajoo/inventory --swagger --chart
```
