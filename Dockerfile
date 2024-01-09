# Builder Image
FROM golang:1.21 as builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . ./

# Turn off CGO since that can result in dynamic links to libc/libmusl which creates problems if you 
# try to run the binary on scratch.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server ./main.go

# Run on secure minimal base image, warning: there's no shell on this!
FROM scratch

COPY --from=builder /bin/server /bin/server

# Doesn't do anything, but it's nice to have so that the engineer running the container knows 
# what port is expected to be published https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8081

CMD ["/bin/server"]