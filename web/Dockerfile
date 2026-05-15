# Build stage
FROM oven/bun:latest AS builder

WORKDIR /app

# Add build argument for the backend URL
ARG VITE_BACKEND_URL
ENV VITE_BACKEND_URL=$VITE_BACKEND_URL

# Copy package files
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Build the application
RUN bun run build

# Final stage
FROM oven/bun:latest

WORKDIR /app

# Copy the build output and necessary files
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

# Expose the port
EXPOSE 3000

# Kamal/kamal-proxy default probe is GET /up; keep Docker health in sync for debugging.
HEALTHCHECK --interval=10s --timeout=5s --start-period=40s --retries=5 \
  CMD bun -e "fetch('http://127.0.0.1:3000/up').then(r=>process.exit(r.ok?0:1)).catch(()=>process.exit(1))"

# TanStack Start's built server is fetch-only (SSR). srvx serves Vite client assets from
# dist/client, then forwards to the handler. --dir + --static must be set this way so paths
# resolve to /app/dist/client (not relative to dist/server/).
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000
CMD ["bun", "node_modules/srvx/bin/srvx.mjs", "--prod", "--host", "0.0.0.0", "--dir", "/app", "--static", "dist/client", "--entry", "dist/server/server.js"]
