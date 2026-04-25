# Quay Go API

[![](https://img.shields.io/badge/go-1.24-version?logo=go&color=rgb(0%2C%20126%2C%20198))]()
[![](https://badgen.net/badge/icon/docker?icon=docker&label)](https://hub.docker.com/r/benjaminfourmauxb/quay-go-api)

An API designed to be deployed alongside the official [Quay Registry](https://github.com/quay/quay) API and to interface with Quay's database to provide features for automation.

## Get stated :rocket:

### Local 💻

Requirements:
- GO `>= 1.24`
- Quay registry setup
  - Database PostgreSQL or MySQL are both supported
  - Quay API token (superuser)


Download dependancies
```
go mod download
```

Run the app
```
export DB_TYPE=postgres # Database type (postrges of mysql)
export DB_DSN="postgres://user:password@localhost:5432/quay" # ConnectionString to the database
export PORT=8080 # (optional) Port where the api is listening (default: 8080)
export LOG_LEVEL=DEBUG # (optional) Adjust the log verbosity. Can be "DEBUG", "INFO", "WARNING" or "ERROR" (default: "DEBUG")
go run main.go
```

Open the swagger, go to [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Generate swagger docs
```
go install github.com/swaggo/swag/cmd/swag@latest
```

```
swag init
```

### Docker image 🐳

Pull the Docker image from [Docker Hub](https://hub.docker.com/r/benjaminfourmauxb/quay-go-api)
```
docker pull benjaminfourmauxb/quay-go-api:latest
```

Run the container
```
docker run -d -e DB_TYPE=postgres -e DB_DSN=postgres://user:password@localhost:5432/quay -p 8080:8080 benjaminfourmauxb/quay-go-api:latest
```

## Features ✨

- **Direct database connection** — Connects directly to the Quay database, bypassing the Quay API layer.
- **Improved error handling** — Structured and consistent error responses for easier debugging and integration.
- **Better REST compliance** — Follows REST principles more closely with proper HTTP methods, status codes, and resource naming.

For more details, see the [Wiki]().

## Version

[![](https://badgen.net/github/tag/BenjaminFourmaux/quay-go-api?cache=600)](https://github.com/BenjaminFourmaux/quay-go-api/tags) [![](https://badgen.net/github/release/BenjaminFourmaux/quay-go-api}?cache=600)](https://github.com/BenjaminFourmaux/quay-go-api/releases)
- [coming soon][v1] First API version with basic actions

## Contributors 👪

[![](https://badgen.net/github/contributors/BenjaminFourmaux/quay-go-api)](https://github.com/BenjaminFourmaux/quay-go-api/graphs/contributors)
- :crown: [Benjamin Fourmaux](https://github.com/BenjaminFourmaux)

## Licence ⚖️

All files on this project is under [**Apache License v2**](https://www.apache.org/licenses/LICENSE-2.0).
You can:
- Reuse the code 
- Modified the code
- Build the code

You must **Mention** the © Copyright if you use and modified code for your own profit. Thank you

© 2026 - Benjamin Fourmaux - All right reserved
