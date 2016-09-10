# klb 

Nash library to mimic the life of [@lborguetti](https://github.com/lborguetti) (aka `klb`).

Ok, just kidding, use `klb` to automate the creation of your infrastructure on AWS.

## Demo

[![asciicast](https://asciinema.org/a/48b1ghj6tli1w0wm3wylqnpyk.png)](https://asciinema.org/a/48b1ghj6tli1w0wm3wylqnpyk)

## Dependencies

- [nash](https://github.com/NeowayLabs/nash)
- [jq](https://stedolan.github.io/jq/)
- [awscli](http://docs.aws.amazon.com/cli/latest/userguide/installing.html)

## Testing

Be careful... Running `make test` will invoke the `aws cli` to create resources on your aws account (if correctly set up)

P.S.: 
- barefoot running is not implemented.
