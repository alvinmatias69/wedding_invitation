FROM golang:1.23 AS build-stage

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /winv cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /winv /winv
COPY ./static /static
COPY ./config.toml /config.toml
COPY ./files /files

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/winv"]
