# Docker Compose Deployment Guide

The [Docker image](https://hub.docker.com/repository/docker/sean1832/wolite) is provided for you to deploy WOLITE on your device with ease.

## 1. Install Docker Container Envrionment

Ubuntu:

```sh
apt install docker.io
```

## 2. Docker Compose Deployment

When using `docker-compose`, the configuration file is as follows:

```yaml
version: "3"
services:
  wolite:
    image: sean1832/wolite:latest
    container_name: wolite
    ports:
      - "3000:3000"
    environment:
      - USERNAME=your-username # required
      - PASSWORD=your-password # required

      # # (Optional) Uncomment to enable
      # - ENABLE_OTP=true
      # - ALLOWED_ORIGINS="ALL" # `ALL` to allow all ip. Add specific ip addresses to restrict access
      # - PORT=3000
    volumes:
      - ./data:/usr/wolite/data
```

Create an empty `data` folder at the current directory where the `docker-compose.yml` file is located.

```sh
mkdir data
```

Run the following command to deploy WOLITE with Docker Compose:

```sh
docker-compose up -d
```

> [!TIP]
> You can find the `OTP_URI` in the `.env` file after the container is started. Use the `OTP_URI` to generate the OTP code with an authenticator app like Google Authenticator or a password manager like 1Password.