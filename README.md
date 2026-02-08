<div align="center">

<img src="https://github.com/sean1832/wolite/blob/main/.design/favicon.svg" alt="wolite logo" width="120" height="120">

# Wolite

![GitHub License](https://img.shields.io/github/license/sean1832/wolite)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org/)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-5.0-FF3E00?logo=svelte)](https://kit.svelte.dev/)
![GitHub Release](https://img.shields.io/github/v/release/sean1832/wolite)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker)](https://hub.docker.com/r/sean1832/wolite)

Wolite is a lightweight, secure Wake-on-LAN (WoL) service that enables remote machine power control over the network.

</div>

## Features

- **Remote Power Control**: Wake up devices on your network with a single click.
- **Secure Authentication**: Protected access with secure login and session management.
- **Device Management**: Add, edit, and manage your devices easily.
- **Single Binary Deployment**: The frontend is embedded directly into the Go binary for easy distribution.
- **Simple Storage**: Uses a local JSON database for simplicity and easy updates.

## Tech Stack

- **Backend**: Go (Golang)
- **Frontend**: Svelte 5 + SvelteKit (TypeScript)
- **Styling**: Tailwind CSS v4 + tailwind-variants
- **UI Components**: shadcn-svelte (bits-ui), Lucide icons
- **Build System**: Taskfile

## Prerequisites

- [Go 1.22+](https://go.dev/)
- [Node.js 20+](https://nodejs.org/)
- [Task](https://taskfile.dev/) (Build tool)

## Docker Deployment

You can deploy Wolite quickly using Docker.

### Option 1: Docker Compose

1.  Create a `data` directory to persist the database and set permission to user 65532:

```bash
mkdir data
chown 65532:65532 data
```

2.  Create a `docker-compose.yml` file:

```yaml
services:
  wolite:
    image: sean1832/wolite:latest
    container_name: wolite
    restart: unless-stopped
    network_mode: host
    volumes:
      # Mount a local directory to persist the database
      # Ensure ./data is writable by user 65532 (or let Docker create it)
      - ./data:/data
    environment:
      # Optional: Override database path within the container
      - DATABASE_PATH=/data/wolite.json
      # Optional: Set the port (default: 8080)
      - PORT=8080
      # Optional: Set JWT secret (if not set, a random one is generated on startup)
      # - JWT_SECRET=your-secure-random-string
      # Optional: Set JWT expiry in seconds (default: 7 days)
      # - JWT_EXPIRY_SECONDS=604800
      # Optional: Enable development mode (allows CORS)
      # - DEV_MODE=false
    user: "65532:65532"
    read_only: true
    tmpfs:
      - /tmp
```

3.  Run the container:

```bash
docker-compose up -d
```

### Option 2: Docker Run

Alternatively, you can run the container directly without `docker-compose`:

```bash
docker run -d \
  --name wolite \
  --restart unless-stopped \
  -p 8080:8080 \
  -v "$(pwd)/data:/data" \
  -e DATABASE_PATH=/data/wolite.json \
  --user "65532:65532" \
  --read-only \
  --tmpfs /tmp \
  sean1832/wolite:latest
```

The application will be available at `http://localhost:8080`.

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd wolite
```

### 2. Build the project

To build both the frontend and backend:

```bash
task build
```

This commands will:

1.  Install frontend dependencies and build the SvelteKit app.
2.  Embed the build artifacts into the Go backend.
3.  Compile the Go binary to `backend/bin/wolite` (or `wolite.exe`).

### 3. Run the application

You need to set up the environment configuration. Create a `.env` file in the `backend` directory or set the environment variables directly.

**Required Environment Variables:**

- `DATABASE_PATH`: Absolute or relative path to the JSON database file (e.g., `./wolite.db` or `/data/wolite.json`).

**Optional Environment Variables:**

- `JWT_SECRET`: Secret key for signing JWTs. One will be generated automatically if not provided.
- `JWT_EXPIRY_SECONDS`: Session expiry time in seconds (default: 604800 / 7 days).
- `DEV_MODE`: Set to `true` to enable CORS (for development).

**Run command:**

```bash
cd backend
./bin/wolite
```

Access the application at `http://localhost:8080`.

## Development

You can run the frontend and backend independently for development features like Hot Module Replacement (HMR).

### Frontend (SvelteKit)

```bash
cd frontend
npm install
npm run dev
```

The frontend will run on `http://localhost:5173`. Update `vite.config.ts` to proxy requests to the backend if needed, or use `DEV_MODE=true` on the backend.

### Backend (Go)

```bash
cd backend
go run main.go
```

## Project Structure

- `backend/`: Go source code, API handlers, and database logic.
- `frontend/`: SvelteKit application, UI components, and styles.
- `Taskfile.yml`: Build scripts and task commands.

## License

[APACHE 2.0](LICENSE)
