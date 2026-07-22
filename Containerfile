# Stage 1: Build SvelteKit frontend
FROM node:22-alpine AS frontend
RUN corepack enable && corepack prepare pnpm@11 --activate
WORKDIR /app
# Vite resolves $env/static/public at build time. The runtime value comes
# from the orchestrator's .env, so any non-empty default works here.
ENV PUBLIC_PB_PORT=8090
ARG VERSION=dev
ENV PUBLIC_APP_VERSION=${VERSION}
COPY sveltekit/package.json sveltekit/pnpm-lock.yaml sveltekit/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY sveltekit/ ./
RUN pnpm build

# Stage 2: Build Go backend
FROM golang:1.25-alpine AS backend
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
# Migrations are compiled into the binary (main.go blank-imports this package)
# and define the schema — the build fails without them. See docs/MIGRATIONS.md.
COPY migrations/ ./migrations/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X github.com/youruser/yourproject/internal/version.Version=${VERSION} -X github.com/youruser/yourproject/internal/version.Commit=${COMMIT} -X github.com/youruser/yourproject/internal/version.Date=${DATE}" -o /server ./cmd/server

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
