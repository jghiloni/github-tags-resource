ARG build_image
FROM ${build_image} AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -trimpath -o /usr/bin/github-tags-resource -ldflags \
    "-w -s -X main.Version=$(cat ./version/version) -X main.BuildRef=$(cat ./.git/short_ref)" .

RUN /usr/bin/github-tags-resource -v

ARG run_image
FROM ${run_image}

USER root

WORKDIR /opt/resource

COPY --from=build /usr/bin/github-tags-resource /usr/bin/github-tags-resource
RUN /usr/bin/github-tags-resource install
