# Docker Deloyment Guide

The [Docker image](https://hub.docker.com/repository/docker/sean1832/wolite) is provided for you to deploy WOLITE on your device with ease.

## 1. Install Docker Container Envrionment

Ubuntu:

```sh
apt install docker.io
```

## 2. Docker Command Deployment

Use the following command to deploy WOLITE with Docker:

```sh
docker run -d -p 3000:3000 \
 -e USERNAME="your-username" \
 -e PASSWORD="your-password" \
 --name wolite
 sean1832/wolite:latest
```

> [!WARNING]
> This command only allows localhost access by default. To allow access from other devices, see [Allow Access from Other Devices](#21-allow-access-from-other-devices) for instruction.

### 2.1 Allow Access from Other Devices

Add the `ALLOWED_ORIGINS` environment variable to allow access from other devices. Use comma to separate multiple IP addresses.

```sh
docker run -d -p 3000:3000 \
  -e USERNAME="your-username" \
  -e PASSWORD="your-password" \
  -e ALLOWED_IPS="ALL" \
  --name wolite
  sean1832/wolite:latest
```

> [!TIP]
> Replace `192.168.x.x` with the IP addresses of the devices you want to allow access from. This is for security reasons to prevent unauthorized access.

### 2.2 Enable One-Time-Password (OTP)

To enable OTP, create an empty `data` folder at the root of the project first, then run the following command with `ENABLE_OTP=true` and mount the `data` directory to the container:

```sh
docker run -d -p 3000:3000 \
 -e USERNAME="your-username" \
 -e PASSWORD="your-password \
 -e ENABLE_OTP=true \
 -v ./data:/usr/wolite/data \
 --name wolite
 sean1832/wolite:latest
```

> [!TIP]
> You can find the `OTP_URI` in the `data/config.json` file after the container is started. Use the `OTP_URI` to generate the OTP code with an authenticator app like Google Authenticator or a password manager like 1Password.
