# WOLITE

![MOCKUP](/docs/images/WOLITE_mockup.png)

A minimal Wake-on-LAN server with a web interface built with **Alpine.js** (frontend) and **Node.js** (backend).

## Installation

### Get Started

```sh
# Clone the repository
git clone https://github.com/sean1832/wolite.git && cd wolite

# Install dependencies
npm install

# Generate a password hash
node generateHash.js
```

Copy the generated hash.

```sh
# Rename the example environment file
mv .env.example .env
```

Edit `.env` and set:

```env
AUTH_PASSWORD_HASH=<your-generated-hash>
```

```sh
# Start the server
node app.js
```

Access the web interface at:  
👉 `http://localhost:3000`

---

### Run in Production

```sh
# Install pm2 globally
npm install -g pm2

# Start the server with pm2
pm2 start app.js
```

Access the web interface at:  
👉 `http://localhost:3000`

#### PM2 Commands

```sh
# Start the process
pm2 start app.js

# Stop the process
pm2 stop app.js

# Restart the process
pm2 restart app.js

# View logs
pm2 logs app.js

# List all processes
pm2 list
```

📄See PM2 doc [here](https://pm2.keymetrics.io/docs/usage/process-management/)

## Tech Stack

- **Node.js**
- **Alpine.js**
- **Tailwind CSS**
- **Express.js**
