FROM golang:latest as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:latest
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

CMD ["sh", "-c", "if [ \"$RUN_TESTS\" = \"true\" ]; then go test -v ./integration-test/... && exit 0; else echo 'Skipping tests'; fi"]