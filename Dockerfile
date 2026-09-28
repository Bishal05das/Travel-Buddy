# FROM golang:1.24 AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download 

# COPY . .

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
#     go build -o server ./cmd/api


# FROM gcr.io/distroless/static-debian12


# COPY --from=builder /app/server .
# COPY migrations /migrations 

# USER nonroot

# EXPOSE 3000

# ENTRYPOINT ["./server"]

FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

RUN mkdir -p /app/images/agencies /app/images/tours

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/server /app/server
COPY --from=builder --chown=nonroot:nonroot /app/images /app/images
COPY --chown=nonroot:nonroot migrations /migrations

ENV MIGRATIONS_PATH=file:///migrations

USER nonroot:nonroot

EXPOSE 3000

ENTRYPOINT ["/app/server"]