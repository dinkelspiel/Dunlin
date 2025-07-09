backend:
    watchexec --restart -- 'go build -o main . && ./main'

frontend:
    cd frontend && pnpm dev