FROM golang:1.21 as build
ENV CGO_ENABLED=0
ENV GOOS=linux
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg
RUN go build -o /bin/go-links ./cmd/go-links

FROM gcr.io/distroless/static-debian12:latest
WORKDIR /
ENTRYPOINT ["/bin/go-links"]
CMD ["-c", "/etc/go-links.yaml"]

COPY --from=build /bin//go-links /bin/go-links
