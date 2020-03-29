FROM scratch
LABEL maintainer="sean@lingrino.com"
COPY glen /
ENTRYPOINT ["/glen"]
