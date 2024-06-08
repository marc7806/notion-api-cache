# syntax=docker/dockerfile:1

FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.22.4-alpine AS build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . ./

RUN go mod download
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /notion-api-cache

FROM alpine:latest

ARG USER=appUser
RUN apk add --update sudo
RUN adduser -D $USER \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER
USER $USER
WORKDIR /home/$USER

COPY --from=build /notion-api-cache ./notion-api-cache

EXPOSE 8080

CMD [ "./notion-api-cache" ]
