FROM golang:1.22 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app ./cmd/main.go

FROM gcr.io/distroless/static-debian11

ENV DB_URL mongodb://db:27017

COPY --from=build /go/bin/app /
CMD ["/app"]
