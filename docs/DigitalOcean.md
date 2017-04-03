Digital Ocean
==

## Authenticating with DigitalOcean

Before proceeding, you'll need to retrieve a DigitalOcean access token to be stored in the doctl configuration file.

This can be done by visiting the Applications & API section of the Control Panel.

You can learn how to generate a token by following the DigitalOcean API guide.

Once you have generated a token, return to your terminal.

To set up doctl and authorize it to use your account, type:

```
doctl auth init
```

You will be prompted to enter the DigitalOcean access token that you generated in the DigitalOcean control panel:

Output
```
DigitalOcean access token: your_DO_token
```

After entering your token, you should receive confirmation that the credentials were accepted:

Output
```
Validating token: OK
```

This will create the necessary directory structure and configuration file to store your credentials.

On OS X and Linux, the configuration file can be found at ${XDG_CONFIG_HOME}/doctl/config.yaml if
the ${XDG_CONFIG_HOME} environmental variable is set, otherwise the config will be written
to ~/.config/doctl/config.yaml. For Windows users, the config will be available at %APPDATA%/doctl/config/config.yaml.
