FROM alpine:3.4

RUN apk update && apk upgrade && apk add ca-certificates

ADD operator /bin/operator

ARG BUILD_DATE
ARG VERSION
ARG VCS_URL
ARG VCS_REF

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="Annotator" \
      org.label-schema.description="Kubernetes operator that adds docker labels as annotations to pods." \
      org.label-schema.vcs-url=$VCS_URL \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.version=$VERSION \
      com.microscaling.license="Apache-2.0" \
      org.label-schema.schema-version="1.0-rc1" 

ENTRYPOINT ["/bin/operator"]
