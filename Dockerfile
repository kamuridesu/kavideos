FROM golang:1.25 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download
RUN PWGO_VER=$(grep -oE "playwright-go v\S+" /modules/go.mod | sed 's/playwright-go //g') \
    && go install github.com/playwright-community/playwright-go/cmd/playwright@${PWGO_VER}

FROM golang:1.25 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /workdir
WORKDIR /workdir
ARG MODULE=bot
RUN GOOS=linux GOARCH=amd64 go build -o /bin/myapp ./${MODULE}/main.go

FROM ubuntu:noble-20251013
EXPOSE 8080
LABEL org.opencontainers.image.authors="kamuridesu@proton.me"

COPY --from=modules /go/bin/playwright /
RUN apt-get update && apt-get install -y ca-certificates tzdata \
    && /playwright install --with-deps \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /bin/myapp /
CMD ["/myapp"]
