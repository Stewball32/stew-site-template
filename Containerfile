# Stage 1: Build SvelteKit frontend
FROM node:22-alpine AS frontend
RUN corepack enable && corepack prepare pnpm@latest --activate
WORKDIR /app
# Vite resolves $env/static/public at build time. The runtime value comes
# from the orchestrator's .env, so any non-empty default works here.
ENV PUBLIC_PB_PORT=8090
COPY sveltekit/package.json sveltekit/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY sveltekit/ ./
RUN pnpm build

# Stage 2: Build Go backend
FROM golang:1.25-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Stage 3: Runtime
FROM alpine:latest
RUN apk add --no-cache ca-certificates wget && \
    adduser -D -u 1000 app
WORKDIR /app
COPY --from=backend /server ./server
COPY --from=frontend /pb_public ./pb_public/
RUN mkdir -p /app/pb_data && chown -R app:app /app
USER app
EXPOSE 8090
VOLUME /app/pb_data
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8090/api/health || exit 1
CMD ["./server", "serve", "--http=0.0.0.0:8090"]
