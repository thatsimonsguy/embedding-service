# syntax=docker/dockerfile:1.4

########## Stage 1: Go builder ##########
FROM golang:1.22-alpine AS go-builder

WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o embedding-service ./main.go


########## Stage 2: llama.cpp builder (static) ##########
FROM alpine:latest AS llama-builder

WORKDIR /llama.cpp
RUN apk add --no-cache build-base git cmake curl-dev linux-headers

RUN git clone https://github.com/ggerganov/llama.cpp . \
    && cmake -S . -B build -DBUILD_SHARED_LIBS=OFF \
    && cmake --build build --target llama-embedding -j$(nproc)

RUN cp ./build/bin/llama-embedding /llama-embedding


########## Stage 3: final container ##########
FROM alpine:latest

WORKDIR /app

# Add runtime dependencies
RUN apk add --no-cache \
    bash \
    libstdc++ \
    libgcc \
    libcurl \
    libgomp

# Copy Go binary
COPY --from=go-builder /app/embedding-service .

# Copy llama embedding binary
COPY --from=llama-builder /llama-embedding ./llama-embedding

# Copy model + optional hash check
COPY models/nomic-embed-text-v1.5.Q4_K_M.gguf /models/nomic-embed-text-v1.5.Q4_K_M.gguf

# ENV config
ENV MODEL_PATH=/models/nomic-embed-text-v1.5.Q4_K_M.gguf \
    EMBEDDING_BINARY=/app/llama-embedding \
    BATCH_SIZE=4096

EXPOSE 8080
ENTRYPOINT ["./embedding-service"]
