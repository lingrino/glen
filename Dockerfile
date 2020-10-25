FROM scratch
LABEL maintainer="sean@lingrino.com"
COPY glen /
ENTRYPOINT ["/glen"]

FROM scratch

# https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL org.opencontainers.image.ref.name="glen" \
    org.opencontainers.image.ref.title="glen" \
    org.opencontainers.image.description="A CLI to gather GitLab project and group variables" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.authors="sean@lingrino.com" \
    org.opencontainers.image.url="https://lingrino.com" \
    org.opencontainers.image.documentation="https://lingrino.com" \
    org.opencontainers.image.source="https://github.com/lingrino/glen"

COPY glen /
ENTRYPOINT ["/glen"]
