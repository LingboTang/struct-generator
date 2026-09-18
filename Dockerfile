FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api-server ./cmd/api-server

FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser
COPY --from=build /out/api-server /usr/local/bin/api-server
USER appuser
EXPOSE 8080
ENTRYPOINT ["api-server"]
