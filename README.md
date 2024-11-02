# Kisara API

**Kisara API** is a RESTful API built with Go and the Fiber framework. This project includes automated configuration, environment variable management, and containerization with Docker to simplify development and deployment.

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Project](#running-the-project)
- [Folder Structure](#folder-structure)
- [License](#license)

## Requirements

- [Go](https://golang.org/doc/install) (latest version recommended)
- [Docker](https://docs.docker.com/get-docker/) (optional, for containerization)
- [Air](https://github.com/cosmtrek/air) (for live reloading during development)

## Installation

Clone this repository to your local directory:

```bash
git clone https://github.com/AnggaaIs/kisara-api-v2.git
cd kisara-api
```

Install the required dependencies:

```bash
go mod download
```

## Configuration

1. Copy `.env.example` to `.env`:

   ```bash
   cp .env.example .env
   ```

2. Adjust the values in `.env` according to your environment setup.

## Running the Project

### Locally

If you have [Air](https://github.com/cosmtrek/air) installed, you can run the server with hot-reloading:

```bash
air
```

To run the server without hot-reloading, use:

```bash
go run ./src/main.go
```

### Using Docker

1. Ensure Docker is installed and running on your system.
2. Build the Docker image with the following command:

   ```bash
   docker build -t kisara-api .
   ```

3. Run the container:

   ```bash
   docker run -v $(pwd)/.env:/app/.env -p 3000:3000 kisara-api
   ```

## Folder Structure

```plaintext
kisara-api
├── src/                   # Main source code
├── tmp/                   # Temporary or cache files
├── .air.toml              # Configuration for Air (hot-reload)
├── .env                   # Environment variables file
├── .env.example           # Example environment variables file
├── .gitignore             # Git ignore file
├── Dockerfile             # Docker configuration to build the image
├── go.mod                 # Go module file
├── go.sum                 # Dependency checksum file
└── LICENSE                # License information
```

## License

This project is licensed under the [GNU Affero General Public License v3.0](./LICENSE). Please refer to the license file for more details.