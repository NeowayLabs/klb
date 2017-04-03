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

## Docs

* [Microsoft Azure](docs/Azure.md)
* [Amazon Web Services](docs/Aws.md)

P.S.:
- barefoot running is not implemented.
