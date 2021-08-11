FROM golang:latest AS go-builder
WORKDIR /build
COPY cmd cmd
COPY internal internal
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go get github.com/ahmetb/govvv
RUN CGO_ENABLED=0 GOOS=linux govvv build -pkg github.com/psidex/CrowsNest/cmd -o ./crowsnest ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=go-builder /build/crowsnest .
ENTRYPOINT ["crowsnest"]
