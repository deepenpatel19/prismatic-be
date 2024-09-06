FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git build-base curl tar pkgconfig
WORKDIR /go/src/app
COPY . .
RUN go mod tidy
RUN export CGO_CFLAGS_ALLOW='-Xpreprocessor'
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/web-app ./main.go

FROM alpine:3.20
RUN export CGO_CFLAGS_ALLOW='-Xpreprocessor'
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
COPY ./migrations/ /usr/bin/migrations/
EXPOSE 8000
ENTRYPOINT /go/bin/web-app