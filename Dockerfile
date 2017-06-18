FROM neowaylabs/klbdeps:0.1

COPY ./aws ${NASHPATH}/lib/klb/aws
COPY ./azure ${NASHPATH}/lib/klb/azure
