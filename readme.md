# WOLITE

![MOCKUP](/docs/images/WOLITE_mockup.png)

A minimal Wake-on-LAN server with a web interface. This project is intended to be a simple and intuitive interface for users to wake up their devices remotely over the same network.

You can deploy this project on a Raspberry Pi or any other device that supports Node.js or Docker.

## Get Started

### Run Locally

```sh
# Install dependencies
npm install

# Set up the environment variables (.env file)
npm run setup

# Start the server
npm run start
```

Access the web interface at:  
👉 `http://localhost:3000`

---

### Run in Production

```sh
# Install pm2 globally
npm install -g pm2

# Start the server with pm2
pm2 start app.js --name wolite
```

Access the web interface at:  
👉 `http://localhost:3000`

---

### Deploy with Docker

Create an empty `.env` file at the current directory:

```sh
touch .env
```

Deploy with docker command:

```sh
docker run -d -p 3000:3000 \
  -e USERNAME="your-username" \
  -e PASSWORD="your-password" \
  -e ALLOWED_IPS="ANY" \
  -e ENABLE_OTP=false \
  -v /full/path/to/.env:/usr/wolite/.env
  --name wolite
  sean1832/wolite
```

For more information, see [Docker Deployment Guide](/docs/deploy-with-docker.md).

> [!TIP]
> Replace `ALLOWED_ORIGINS="ANY"` with specific IP addresses to restrict access to the web interface. `ANY` allows all origins. You can use a comma-separated list of IP addresses to allow multiple origins (e.g., `ALLOWED_IPS="192.168.x.x, 192.168.x.x"`).

> [!TIP]
> You must provide **FULL PATH** to the `.env` file in the `-v` flag. Replace `/full/path/to/.env` with the actual path to the `.env` file.

### Deploy with Docker Compose

Create an empty `.env` file at the current directory where the `docker-compose.yml` file is located.

```sh
touch .env
touch docker-compose.yaml
```

Copy and paste the following to `docker-compose.yaml`

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
      # - ENABLE_OTP=true
      # - ALLOWED_ORIGINS="ANY" # wildcard `ANY` to allow all ip. Add specific ip addresses to restrict access
    volumes:
      - ./.env:/usr/wolite/.env
```

Run the following command to deploy WOLITE with Docker Compose:

```sh
docker-compose up -d
```

> [!TIP]
> You can find the `OTP_URI` in the `.env` file after the container is started. Use the `OTP_URI` to generate the OTP code with an authenticator app like Google Authenticator or a password manager like 1Password.

For more information, see [Docker Compose Deployment Guide](/docs/deploy-with-docker-compose.md).
