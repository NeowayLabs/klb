# Microsoft Azure

## Service principal and Setting account

The first thing you need to set up is your credentials, that are also
used to define where the infrastructure changes will be applied.

We have a tools to make this easier, its output can be used
to set environment variables on a easy way for you.

It is a good idea to save this env vars on a safe place on your host
to make your life easier.

### Creating a service principal

Before anything, you must be logged before creating a service principal,
so run:

```
azure login
```

Now that you are logged, first you need to
define in which subscription you want to create it.

To know which subscriptions are available for your user run:

```sh
azure account list
```

You will need the **ID** of the subscription for the next step.
When you have the subscription **ID** just run:

```sh
./tools/azure/createsp.sh <subscription-id> <service-principal-name> <password>
```

If it executes with success you will have a valid service principal
to start having some fun at Azure.

### Setting account

If you already have a `service principal` configured, you can setup
the environment variables using the `getcredentials.sh` script.

You'll need the environment variables below:

- AZURE_SUBSCRIPTION_ID=&lt;subscription id&gt;
- AZURE_TENANT_ID=&lt;tenant id&gt;
- AZURE_CLIENT_ID=&lt;AppId of service principal&gt;
- AZURE_CLIENT_SECRET=&lt;password of service principal&gt;


Basic usage:

```sh
λ> ./tools/azure/getcredentials.sh
Usage:  ./tools/azure/getcredentials.sh <(sh|nash)> <subscription name> <service principal name> <service secret>
```

For cool nash shells:

```sh
λ> ./tools/azure/getcredentials.sh nash subscription-name klb-sp-tests 123456
setenv AZURE_SUBSCRIPTION_ID="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXXXXX"
setenv AZURE_SUBSCRIPTION_NAME="CLOUDFORMATION - QA"
setenv AZURE_TENANT_ID="XXXXXXXX-XXXX-XXXX-XXXXXXXXXXXXXX"
setenv AZURE_CLIENT_ID="XXXXXXXX-XXXX-XXXX-XXXXXXXXXXXXXX"
setenv AZURE_CLIENT_SECRET="123456"
```

Copy the environment variables from the output to a file and import into your nash session:

```sh
λ> import ./credentials
```

For lame linux shell (like bash):

```sh
λ> ./tools/azure/getcredentials.sh sh subscription-name klb-sp-tests 123456
export AZURE_SUBSCRIPTION_ID="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXXXXX"
export AZURE_SUBSCRIPTION_NAME="CLOUDFORMATION - QA"
export AZURE_TENANT_ID="XXXXXXXX-XXXX-XXXX-XXXXXXXXXXXXXX"
export AZURE_CLIENT_ID="XXXXXXXX-XXXX-XXXX-XXXXXXXXXXXXXX"
export AZURE_CLIENT_SECRET="123456"
```

Write the output above to a file and import into your lame shell session:

```sh
λ> source ./credentials
```

Besides the credential your subscription must also be registered
on some resource provider namespaces, since we will be creating resources.

To register on the required providers you can run:

```sh
λ> ./tools/azure/registerproviders.sh
```

## Examples

There are some examples on the **examples/azure** dir. To build them
you can call:

```
make example-<name>
```

Where **<name>** can be the name of any example on the dir.

Running the backup example:

```
make example-backup
```
