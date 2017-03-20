# klb

[![Build Status](https://travis-ci.org/NeowayLabs/klb.svg?branch=master)](https://travis-ci.org/NeowayLabs/klb)

Nash library to mimic the life of [@lborguetti](https://github.com/lborguetti) (aka `klb`).

Ok, just kidding, use `klb` to automate the creation of your
infrastructure on AWS or Azure.

## Demo

[![asciicast](https://asciinema.org/a/48b1ghj6tli1w0wm3wylqnpyk.png)](https://asciinema.org/a/48b1ghj6tli1w0wm3wylqnpyk?autoplay=true&speed=2)

## Dependencies

- [nash](https://github.com/NeowayLabs/nash)
- [jq](https://stedolan.github.io/jq/)
- [awscli](http://docs.aws.amazon.com/cli/latest/userguide/installing.html)

You can run the following command to install deps:

```sh
make deps
```

## Install

Just run:

```sh
make install
```

To install klb on your **NASHPATH**.


### Updating vendored dependencies

```sh
make vendor
```

## Testing

Just run `make testall`.

For each cloud you'll need the environment variables.
See the docs for each cloud to help you with each one.

Logging by default will be saved on files, since the tests can be
pretty long running and you can check out the progress on the files.

Inside each test package the logs will be saved at **./testdata/logs**.

To run redirecting logs to stdout:

Just run `make testall logger=stdout`.


### Azure

You'll need the environment variables below:

- AZURE_SUBSCRIPTION_ID=&lt;subscription id&gt;
- AZURE_TENANT_ID=&lt;tenant id&gt;
- AZURE_CLIENT_ID=&lt;AppId of service principal&gt;
- AZURE_CLIENT_SECRET=&lt;password of service principal&gt;

If you already have a `service principal` configured, you can setup
the environment variables using the `getcredentials.sh` script.

Basic usage:

```sh
λ> ./tools/azure/getcredentials.sh
Usage:  ./tools/azure/getcredentials.sh <(sh|nash)> <service principal name> <service secret>
```

For cool nash shells:

```sh
λ> ./tools/azure/getcredentials.sh nash klb-sp-tests 123456
setenv AZURE_SUBSCRIPTION_ID="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXXXXX"
setenv AZURE_SUBSCRIPTION_NAME="CLOUDFORMATION - QA"
setenv AZURE_TENANT_ID="XXXXXXXX-XXXX-XXXX-XXXXXXXXXXXXXX"
setenv AZURE_CLIENT_ID="XXXXXXXX-XXXX-XXXX-XXXXXXXXXXXXXX"
setenv AZURE_CLIENT_SECRET="123456"
```

Redirect the output above to a file and import into your nash session:

```sh
λ> ./tools/azure/getcredentials.sh nash <service principal name> <service secret> > credentials
λ> import ./credentials
```

For lame linux shell (like bash):

```sh
λ> ./tools/azure/getcredentials.sh sh klb-sp-tests 123456
export AZURE_SUBSCRIPTION_ID="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXXXXX"
export AZURE_SUBSCRIPTION_NAME="CLOUDFORMATION - QA"
export AZURE_TENANT_ID="XXXXXXXX-XXXX-XXXX-XXXXXXXXXXXXXX"
export AZURE_CLIENT_ID="XXXXXXXX-XXXX-XXXX-XXXXXXXXXXXXXX"
export AZURE_CLIENT_SECRET="123456"
```

Redirect the output above to a file and import into your lame shell session:

```sh
λ> ./tools/azure/getcredentials.sh sh <service principal name> <service secret> > credentials
λ> source ./credentials
```

Besides the credential your subscription must also be registered
on some resource provider namespaces, since we will be creating resources.

To register on the required providers you can run:

```sh
λ> ./tools/azure/registerproviders.sh
```

If you have not configured the service principal yet, then the
following section could be helpful.


#### Creating a service principal

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

P.S.:
- barefoot running is not implemented.
