FROM golang:1.24.5-alpine AS builder

RUN apk add --no-cache gcc g++ make pkgconfig

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o govertor ./cmd/govertor/main.go


FROM alpine:latest

RUN apk add --no-cache ffmpeg fontconfig ttf-dejavu

WORKDIR /app

COPY --from=builder /app/govertor .

COPY assets ./assets
COPY converted ./converted

ENTRYPOINT ["./govertor"]

