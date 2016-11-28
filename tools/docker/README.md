Docker image with [nash][1] and [klb][2]
===

# Create Image

        docker -t neowaylabs/klb .

## Run with AWS

        docker run --rm -ti \
            -e "AWS_ACCESS_KEY_ID=xxxxxxxxxx" \
            -e "AWS_SECRET_ACCESS_KEY=xxxxxx" \
            -e "AWS_DEFAULT_REGION=us-west-2"
            neowaylabs/klb nash

         import klb/aws/all



## Run with Azure

        docker run --rm -ti \
            -e "AZURE_SUBSCRIPTION_ID=xxxxxx" \
            -e "AZURE_SUBSCRIPTION_NAME=xxxx" \
            -e "AZURE_TENANT_ID=xxxxxxxxxxxx" \
            -e "AZURE_CLIENT_ID=xxxxxxxxxxxx" \
            -e "AZURE_CLIENT_SECRET=xxxxxxxx" \
            neowaylabs/klb nash

        import klb/azure/all


[1]:https://github.com/NeowayLabs/nash
[2]:https://github.com/NeowayLabs/klb
