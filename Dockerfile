FROM neowaylabs/klbdeps:0.4

COPY ./aws ${NASHPATH}/lib/klb/aws
COPY ./azure ${NASHPATH}/lib/klb/azure

COPY ./tools/azure/createsp.sh ${NASHPATH}/bin/azure-createsp.sh
COPY ./tools/azure/getcredentials.sh ${NASHPATH}/bin/azure-getcredentials.sh

ENV PATH $PATH:${NASHPATH}/bin
