# klb

[![Build Status](https://travis-ci.org/NeowayLabs/klb.svg?branch=master)](https://travis-ci.org/NeowayLabs/klb)

Nash library to mimic the life of [@lborguetti](https://github.com/lborguetti) (aka `klb`).

Ok, just kidding, use `klb` to automate the creation of your
infrastructure on AWS or Azure.

## Demo

[![asciicast](https://asciinema.org/a/48b1ghj6tli1w0wm3wylqnpyk.png)](https://asciinema.org/a/48b1ghj6tli1w0wm3wylqnpyk?autoplay=true&speed=2)

## Dependencies

- python pip
- nodejs npm
- [nash](https://github.com/NeowayLabs/nash)
- [jq](https://stedolan.github.io/jq/)
- [awscli](http://docs.aws.amazon.com/cli/latest/userguide/installing.html)
- [azure-cli 1.0](https://github.com/Azure/azure-xplat-cli)
- azure-cli 2.0

To aid you we provide some scripts to install the deps, per distro.
For arch linux you can use:

```sh
sudo ./deps/arch.sh
```

Or for debian linux you can use:

```sh
sudo ./deps/debian.sh
```

Contributions for other distros are welcomed.
The scripts can assume that Go is already installed, all the
other dependencies will be installed automatically.

## Install

Just run:

```sh
make install
```

To install klb on your **NASHPATH**.

## Running with Docker

A **neowaylabs/klb** image is also provided with the project,
saving you the hussle of installing the dependencies on
your machine.

To create a fast nash shell where you can play around with
klb you can run:

```
make shell
```

To enter the shell you need to export the required variables
that will enable you to use klb. More details on getting
your credentials and exporting them on your environment
can be found bellow.

The image is ready to run nash scripts that import klb modules.

### Updating vendored dependencies

```sh
make vendor
```

## Testing

Just run `make test`.

For each cloud you'll need the environment variables.
See the docs for each cloud to help you with each one.

Logging by default will be saved on files, since the tests can be
pretty long running and you can check out the progress on the files.

Inside each test package the logs will be saved at **./testdata/logs**.

To run redirecting logs to stdout:

Just run `make test logger=stdout`.

There are also examples that can be run automatically, to validate
if they are working. Just run:

```
make test-examples
```

They are not included on the CI or the common tests because they take
too much time to run, but it is a good way to validate that complex
scenarios are working fine.

## Docs

* [Microsoft Azure](docs/Azure.md)
* [Amazon Web Services](docs/Aws.md)

P.S.:
- barefoot running is not implemented.

## Cleanup

The automated tests strives to always cleanup ALL created resources.
But sometimes it may fail to delete resources, it can happen even
because of a intermitent cloud service failure.

If you want to be absolutely sure to delete all test resources
created by klb run:

```
make cleanup
```

Do not worry, resources are just going to be deleted
after you carefully check the list of resources and
accept it, it won't go beserk deleting everything on
your account.
