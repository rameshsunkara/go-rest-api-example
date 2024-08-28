FROM golang:1.22 as builder
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
ENTRYPOINT ["./app"]
EXPOSE $port
