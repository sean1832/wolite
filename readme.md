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
pm2 start app.js
```

Access the web interface at:  
👉 `http://localhost:3000`

#### PM2 Commands

```sh
# Start the server with pm2 and name the process
pm2 start app.js --name wolite

pm2 save # Save the process list
pm2 startup # set up the startup script
```

📄See PM2 doc [here](https://pm2.keymetrics.io/docs/usage/process-management/)

## Tech Stack

- **Node.js**
- **Alpine.js**
- **Tailwind CSS**
- **Express.js**
