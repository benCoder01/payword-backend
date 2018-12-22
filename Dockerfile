# api Dockerfile

FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR /go/src/gitlab.com/benCoder01/payword-backend
COPY . .

RUN dep ensure --vendor-only

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/payword-backend

#RUN go build -o /go/bin/payword-backend main.go

#CMD ["./main"]
#CMD [ "go run main.go" ]

FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/payword-backend /go/bin/payword-backend

ENTRYPOINT ["/go/bin/payword-backend"]