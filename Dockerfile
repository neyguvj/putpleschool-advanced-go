FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/migrations ./cmd/migrations \
 && CGO_ENABLED=0 go build -o /bin/server ./cmd/server

FROM alpine:3.22

COPY --from=build /bin/migrations /bin/server /bin/

CMD ["/bin/server"]
