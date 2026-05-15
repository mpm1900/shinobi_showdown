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
# TanStack Start / Vinxi usually puts the server in .output
COPY --from=builder /app/.output ./.output
COPY --from=builder /app/package.json ./package.json

# Expose the port (TanStack Start defaults to 3000)
EXPOSE 3000

# Start the server
CMD ["bun", "run", ".output/server/index.mjs"]
