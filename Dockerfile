FROM golang:1.20 as builder
LABEL stage=builder
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies 
# and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy source files and build the binary
COPY . .
RUN make build


FROM scratch
WORKDIR /app/
ARG port
COPY --from=builder /usr/src/app/app .
COPY --from=builder /usr/src/app/config/*.yaml /app/config/
## As of golang:1.20 this cerrtificate is included in the base image.
## So, it is not needed to add this certificate manually or run any add ca-certificates commands.
## This certificate is added to avoid any issues incase your application is making any HTTPS etc., connections
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./app"]
EXPOSE $port
