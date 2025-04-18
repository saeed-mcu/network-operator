# network-operator

## Technical requirements
* go version 1.16+
* An operator-sdk binary installed locally
* Make sure your user is authorized with cluster-admin permissions.

## Step 1 : Setting up your project

initialize a boilerplate project structure with the following:

```bash
# we'll use a domain of digicloud.ir
# so all API groups will be <group>.digicloud.ir
operator-sdk init --domain digicloud.ir --repo github.com/saeed-mcu/network-operator
```

## Step 2 : Defining an API
```bash
operator-sdk create api --group network --version v1 --kind DigiNet --resource --controller
```
Chnage the code and logic for reconciler and type:
```
make generate
make manifest
```
