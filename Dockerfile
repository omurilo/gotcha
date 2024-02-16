FROM golang:1.22 as build

WORKDIR /go/src/app

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app ./cmd/main.go

FROM gcr.io/distroless/static-debian11

ENV DB_URL mongodb://db:27017
ENV GOGC 250

COPY --from=build /go/bin/app /
CMD ["/app"]
