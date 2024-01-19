FROM golang:1.21 as build
ENV CGO_ENABLED=0
ENV GOOS=linux
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY pkg pkg
RUN go build -o /go-links .

FROM scratch
WORKDIR /
ENTRYPOINT ["/bin/go-links", "-c", "/etc/go-links.yaml"]

COPY --from=build /go-links /bin/go-links
COPY go.yaml /etc/go-links.yaml
