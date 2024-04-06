FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-socket

FROM gcr.io/distroless/base-debian11 AS build-release-state

WORKDIR /

COPY --from=build-stage /go-socket /go-socket

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/go-socket"]