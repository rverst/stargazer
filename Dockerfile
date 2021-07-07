FROM golang:1.16 AS builder

ARG VERSION=unknown

WORKDIR $GOPATH/src/github.com/rverst/stargazer

ENV CGO_ENABLED 0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN go build -ldflags="-X 'main.version=${VERSION}'" -o /stargazer

FROM alpine

COPY --from=builder /stargazer /usr/bin/stargazer
RUN mkdir /data

ENV OUTPUT_FILE "/data/README.md"
ENV GITHUB_USER ""
ENV ACCESS_TOKEN ""
ENV IGNORE_REPOS ""

ENTRYPOINT ["stargazer"]
