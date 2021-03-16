FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-openshift-4.8 AS builder

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /src

# copy just enough of the git repo to parse HEAD, used to record version in OLM binaries
COPY .git/HEAD .git/HEAD
COPY .git/refs/heads/. .git/refs/heads

ARG STAGING_DIR
RUN echo ${STAGING_DIR}
COPY ${STAGING_DIR} .
COPY vendor vendor
COPY Makefile Makefile
RUN make build

# copy and build vendored grpc_health_probe
RUN CGO_ENABLED=0 go build -mod=vendor -tags netgo -ldflags "-w" ./vendor/github.com/grpc-ecosystem/grpc-health-probe/...

FROM registry.ci.openshift.org/ocp/4.8:base

COPY --from=builder /src/bin/* /bin/registry/
COPY --from=builder /src/grpc-health-probe /bin/grpc_health_probe

RUN ln -s /bin/registry/* /bin

RUN mkdir /registry
RUN chgrp -R 0 /registry && \
    chmod -R g+rwx /registry
WORKDIR /registry

# This image doesn't need to run as root user
USER 1001

EXPOSE 50051

ENTRYPOINT ["/bin/registry-server"]
CMD ["--database", "/bundles.db"]

LABEL io.k8s.display-name="OpenShift Operator Registry" \
      io.k8s.description="This is a component of OpenShift Operator Lifecycle Manager and is the base for operator catalog API containers." \
      maintainer="Odin Team <aos-odin@redhat.com>" \
      summary="Operator Registry runs in a Kubernetes or OpenShift cluster to provide operator catalog data to Operator Lifecycle Manager."
