FROM golang:1.22 AS build

WORKDIR /go/src/boilerplate

COPY . .
COPY ./cmd/main.go .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o dpo-test-backend .

FROM alpine:latest as release

WORKDIR /app

COPY --from=build /go/src/boilerplate .
RUN rm -rf ./main.go

EXPOSE 8080/tcp

CMD ["./dpo-test-backend"]