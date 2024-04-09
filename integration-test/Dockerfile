FROM golang:latest as modules
COPY ../go.mod ../go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:latest
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

CMD ["go", "test", "-v", "./integration-test/..."]