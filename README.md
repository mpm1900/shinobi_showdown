# Ninja Battler V1

<img width="1403" height="999" alt="Screenshot 2026-05-02 at 2 19 33 PM" src="https://github.com/user-attachments/assets/47dac756-a643-4548-8395-d3e8ccad8f0c" />
<img width="1518" height="948" alt="Screenshot 2026-05-04 at 11 14 29 PM" src="https://github.com/user-attachments/assets/6dacafbf-7e2c-4762-b040-8d887deff610" />
<img width="1529" height="947" alt="Screenshot 2026-05-04 at 11 15 10 PM" src="https://github.com/user-attachments/assets/192f4a73-cb6e-4098-a68f-3217c7697db2" />


A modern web-based turn-based combat game inspired by the Naruto universe. Features an animated background, real-time combat logic, and a TanStack-powered frontend.

## Prerequisites

Ensure you have the following installed:

- **Go**: 1.25+
- **Docker & Docker Compose**: For the PostgreSQL database
- **Bun**: (Recommended) or Node.js/pnpm for the frontend
- **Goose**: For database migrations (`go install github.com/pressly/goose/v3/cmd/goose@latest`)
- **Air**: (Optional) For Go hot-reloading (`go install github.com/air-verse/air@latest`)

## Getting Started

### 1. Environment Setup

Copy the example environment file and adjust if necessary:

```bash
cp .env.example .env
```

### 2. Infrastructure

Start the PostgreSQL database using Docker Compose:

```bash
make up
```

### 3. Database Migrations

Run the migrations to set up the database schema:

```bash
make migrate
```

### 4. Running the Backend

You can run the backend in two ways:

**Standard run:**
```bash
make run
```

**With hot-reloading (requires Air):**
```bash
air
```

The server will start on the port specified in your `.env` (default is `:3005`).

### 5. Running the Frontend

Navigate to the `web` directory, install dependencies, and start the development server:

```bash
cd web
bun install
bun dev
```

The frontend will be available at [http://localhost:3000](http://localhost:3000).
