# Simeru Scraper

Simeru Scraper is a web scraper that scrapes schedules from the official website and caches the data in Redis. The API is built with Go and Fiber.

[![Go](https://img.shields.io/github/go-mod/go-version/savioruz/simeru-scraper)](https://golang.org/)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/savioruz/simeru-scraper)
[![Go Report Card](https://goreportcard.com/badge/github.com/savioruz/roastgithub-api)](https://goreportcard.com/report/github.com/savioruz/simeru-scraper)
[![GitHub issues](https://img.shields.io/github/issues/savioruz/simeru-scraper)](https://goreportcard.com/report/github.com/savioruz/simeru-scraper)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/savioruz/simeru-scraper)](https://goreportcard.com/report/github.com/savioruz/simeru-scraper)

## Table of Contents

- [Features](#features)
- [Deployment](#deployment)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [Running the API](#running-the-api)
  - [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

## Features

- Scrapes schedules from the official website
- Caches data in Redis
- Cron job to update data
- API Documentation with Swagger
- Docker support

## Deployment

- ### Koyeb
[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/services/deploy?type=git&builder=dockerfile&repository=github.com/savioruz/roastgithub-api&branch=main&ports=3000;http;/&name=simeru-scraper-koyeb&env[APP_HOST]=0.0.0.0&env[APP_PORT]=3000&env[REDIS_HOST]=YOUR_REDIS_HOST&env[REDIS_PORT]=6379&env[REDIS_PASSWORD]=&env[REDIS_DB_NUMBER]=0)

## Requirements

- Go 1.23+
- Docker
- Redis
- Make

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/savioruz/simeru-scraper.git
    cd simeru-scraper
    ```

2. **Environment Variables:**

   Create a `.env` file in the root directory and add the following:

    ```bash
    cp .env.example .env
    ```

## Usage

### Running the API

You can run the API using Docker or directly with Make.

### Docker (Recommended)

1. **Run redis:**

    ```bash
   make docker.redis
   ```

2. **Run the application:**

    ```bash
    make docker.run
    ```

For production, you need to secure redis on Makefile with a password.

### Make

1. **Run the application:**

    ```bash
    make run
    ```

You need to have Redis running on your machine.

### API Documentation

Swagger documentation is available at: http://localhost:3000/swagger.

![Preview](/assets/preview.png)

## Project Structure

```
.
├── config/                 # Configuration files
│── docs/                   # Project documentation
├── internal/
│   ├── adapters/           # Adapters for external services
│   │   ├── cache/          # Cache layer
│   │   ├── handlers/       # Handlers layer
│   │   └── repositories/   # Storage layer integration
│   └── cores/              # Core business layer
│       ├── entities/       # Business entities
│       ├── ports/          # Adapter implementations
│       │   └── ports.go
│       └── services/       # Use cases layer
├── pkg/
│   ├── constant/
│   ├── middleware/
│   ├── routes/
│   ├── server/            # Server configuration
│   └── utils/             # Utility functions
├── main.go
├── .env
└── Dockerfile

```

## Contributing

Feel free to open issues or submit pull requests with improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Fiber](https://github.com/gofiber/fiber)
