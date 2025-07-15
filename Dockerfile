# Build stage
FROM alpine:3.21

RUN apk add --no-cache \
    "go=~1.23" \
    nodejs \
    npm

RUN npm install -g pnpm

WORKDIR /app

COPY frontend ./frontend
WORKDIR /app/frontend
RUN echo Y | pnpm install && pnpm build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
