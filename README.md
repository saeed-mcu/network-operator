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
