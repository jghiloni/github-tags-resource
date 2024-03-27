ARG build_image
ARG run_image
ARG version
ARG build_ref
FROM ${build_image} AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /build

COPY . .

RUN ["go", "mod", "download"]
RUN ["go", "build", "-trimpath", "-o", "/usr/bin/github-tags-resource", "-ldflags",\
    '"-w -s -X main.Version=${version} -X main.BuildRef=${build_ref}"', "."]

RUN ["/usr/bin/github-tags-resource", "-v"]

FROM ${run_image} AS run

USER root

WORKDIR /opt/resource

COPY --from=build /usr/bin/github-tags-resource /usr/bin/github-tags-resource
RUN ["/usr/bin/github-tags-resource", "install"]
