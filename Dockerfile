# syntax=docker/dockerfile:1

FROM golang:1.19 AS build

WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download

COPY ./src/*.go .

RUN GOOS=linux go build -o ./discovery .

# FROM alpine:3.17
# FROM golang:1.19
FROM gcr.io/distroless/base-debian11
# RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=build /app/discovery ./
USER nonroot:nonroot

ENV SLEEP=5
ENV MODE="loop"
# ENV MODE="watch"
ENV DISCOVERY="endpoints"
ENV OUTPUT_FORMATS="default_yaml,traefik_yaml"
ENV LABEL_FILTER="k8s-service-discovery"
ENV MINIO_REGION="us-east-1"

CMD ["./discovery"]
