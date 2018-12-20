# api Dockerfile

FROM golang:1.10

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR /go/src/gitlab.com/benCoder01/payword-backend
COPY . .

RUN dep ensure --vendor-only

RUN go build main.go

CMD ["./main"]
#CMD [ "go run main.go" ]