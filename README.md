# klb

Nash library to mimic the life of [@lborguetti](https://github.com/lborguetti) (aka `klb`).

Ok, just kidding, use `klb` to automate the creation of your
infrastructure on AWS or Azure.

## Demo

[![asciicast](https://asciinema.org/a/48b1ghj6tli1w0wm3wylqnpyk.png)](https://asciinema.org/a/48b1ghj6tli1w0wm3wylqnpyk?autoplay=true&speed=2)

## Dependencies

- [nash](https://github.com/NeowayLabs/nash)
- [jq](https://stedolan.github.io/jq/)
- [awscli](http://docs.aws.amazon.com/cli/latest/userguide/installing.html)

## Testing

To run `make testazure` you'll need the environment variables below:

- AZURE_SUBSCRIPTION_ID=&lt;subscription id&gt;
- AZURE_TENANT_ID=&lt;tenant id&gt;
- AZURE_CLIENT_ID=&lt;AppId of service principal&gt;
- AZURE_CLIENT_SECRET=&lt;password of service principal&gt;

The values of `AZURE_SUBSCRIPTION_ID` and `AZURE_TENANT_ID` could be
obtained from the command below:

```sh
λ> azure account show
info:    Executing command account show
data:    Name                        : <SUBSCRIPTION NAME>
data:    ID                          : <AZURE_SUBSCRIPTION_ID>
data:    State                       : Enabled
data:    Tenant ID                   : <AZURE_TENANT_ID>
data:    Is Default                  : true
data:    Environment                 : AzureCloud
data:    Has Certificate             : No
data:    Has Access Token            : Yes
data:    User name                   : <your email address>
data:
info:    account show command OK

```

In the output above, `ID` is the `AZURE_SUBSCRIPTION_ID` and `Tenant ID`
is the `AZURE_TENANT_ID`. The values of `AZURE_CLIENT_ID` and
`AZURE_CLIENT_SECRET` came from a previously created service principal.
To create a simple service principal to authenticate in the API,
follow the steps below:

The command below create a service principal called `klb-sp-tests`
with a password of your choice. This password will be the value of
`AZURE_CLIENT_SECRET` environment variable.

```sh
λ> azure ad sp create -n klb-sp-tests -p <your password>
```

You'll need to grant permissions to the service principal authenticate
on behalf of your subscription id.
You'll need the Object ID of the just created `service principal`. To
get this value, run the command below and look for the service name
`klb-sp-tests`:

```sh
λ> azure ad sp list
```

And then issue the command below to grant permissions to the `klb-sp-tests`:

```sh
azure role assignment create --objectId <klb-sp-tests object id> -o Reader -c /subscriptions/{AZURE_SUBSCRIPTION_ID}/
```

If everything worked as expected, then export the required environment
variables and run:

```
λ> make testazure
```

P.S.:
- barefoot running is not implemented.
