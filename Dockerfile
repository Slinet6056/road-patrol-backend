FROM golang:1.22 as builder

LABEL org.opencontainers.image.source="https://github.com/Slinet6056/road-patrol-backend"
LABEL org.opencontainers.image.licenses="MIT"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o roadpatrol ./cmd

FROM scratch

COPY --from=builder /app/roadpatrol .

ENTRYPOINT ["./roadpatrol"]
