# network-operator


## Step 1 : Setting up your project

```bash
operator-sdk init --domain digicloud.ir --repo github.com/saeed-mcu/network-operator
```

## Step 2 : Defining an API
```bash
operator-sdk create api --group network --version v1 --kind NetConf --resource --controller
```
Chnage the code and logic for reconciler and type:
```
make generate
make manifest
```

### Vendoring mode for build
```
export GOPROXY=https://repo.hami.digicloud.ir/repository/go-proxy/,direct
go mod vendor
```
