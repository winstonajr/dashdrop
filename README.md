<div align="center">

# ⚡ DashDrop

**Blazing-fast, self-hosted local area network (LAN) file exchange service with zero cloud dependencies.**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Framework](https://img.shields.io/badge/Framework-Gin-008ECF?style=flat-square)](https://gin-gonic.com/)
[![Database](<https://img.shields.io/badge/Database-SQLite%20(sqlc)-003B57?style=flat-square&logo=sqlite>)](https://sqlite.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

[Features](#-features) •
[Architecture](#-architecture) •
[Getting Started](#-getting-started) •
[Network Access](#-accessing-from-other-devices) •
[API Reference](#-api-reference) •
[License](#-license)

</div>

---

## 📌 Overview

**DashDrop** is a lightweight, high-performance self-hosted utility designed to eliminate external cloud transfer limits and physical flash drives across office or home networks.

Built in Go with SQLite persistence via `sqlc`, DashDrop streams multipart uploads directly to disk without bloating system memory, delivering a dark-themed, responsive web client featuring real-time search, drag-and-drop, and resilient link sharing.

---

## ✨ Features

- **Chunked Direct Streaming:** Files are streamed straight to disk buffers via `io.Copy`, keeping memory footprints minimal regardless of file size.
- **Type-Safe SQLite Persistence:** Queries and models generated via `sqlc` with automated SQLite migrations.
- **Collision-Proof Storage:** Files are assigned UUIDv4 identifiers on the filesystem while preserving original metadata and download filenames.
- **Modern Glassmorphic UI:** Built with Vue 3 and Tailwind CSS, featuring file-type icon detection, upload progress tracking, and instant search (press `/`).
- **LAN-Ready Clipboard Support:** Automatic fallback to `document.execCommand` ensures one-click download link copying over unencrypted HTTP local network IPs.
- **Single-Binary Mindset:** Gin serves both REST API endpoints and static assets directly from port `8080`.

---

## 🏗 Architecture

The backend implements a clean, layered architecture:

```text
dashdrop/
├── cmd/
│   └── api/             # Application entrypoint & dependency injection
├── internal/
│   ├── config/          # Environment & server configurations
│   ├── dto/             # Data Transfer Objects
│   ├── handler/         # HTTP controllers & Gin request binding
│   ├── repository/      # sqlc-generated database access layer
│   ├── router/          # Endpoint definitions & static file routing
│   └── service/         # Core business logic (I/O, rollback, disk operations)
├── public/
│   └── index.html       # Single Page Application (Vue 3, Tailwind CSS)
├── sql/
│   ├── migrations/      # DDL schemas for SQLite
│   └── queries.sql      # Type-annotated SQL queries
├── uploads/             # File storage destination (ignored by git)
└── sqlc.yaml            # sqlc compiler configuration
```

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.22+
- A C compiler (such as `gcc`) for `mattn/go-sqlite3` (CGO)
- [Air](https://github.com/air-verse/air) for live-reloading (optional)

### 1. Clone the repository

```bash
git clone https://github.com/winstonajr/dashdrop.git
cd dashdrop
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the server

**Development mode (live-reload with Air):**

```bash
air
```

**Standard production run:**

```bash
CGO_ENABLED=1 go run ./cmd/api
```

The server binds to `0.0.0.0:8080` by default.

---

## 🌐 Accessing from Other Devices

1. Identify your local IP address:
   ```bash
   # On Linux:
   ip -br addr show
   ```
2. Open a web browser on any smartphone, tablet, or secondary computer connected to the same Wi-Fi or LAN:
   ```text
   http://192.168.X.X:8080
   ```

---

## 🔌 API Reference

| Method   | Endpoint            | Description                                       |
| -------- | ------------------- | ------------------------------------------------- |
| `GET`    | `/`                 | Serves the web client SPA                         |
| `GET`    | `/api/files`        | Lists all stored file metadata in JSON            |
| `POST`   | `/api/upload`       | Uploads file via `multipart/form-data`            |
| `GET`    | `/api/download/:id` | Downloads the file with its original filename     |
| `DELETE` | `/api/files/:id`    | Permanently deletes file metadata and disk record |

---

## 📄 License

Distributed under the **MIT License**. See [`LICENSE`](LICENSE) for more information.
