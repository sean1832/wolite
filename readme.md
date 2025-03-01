# WOLITE

A minimal Wake on LAN server with a web interface written with Alpine.js(frontend) and Node.js(backend).

## Installation

1. Clone the repository
2. Run `npm install`
3. Generate a hash of your desired password by running `node generateHash.js` and copy the output.
4. Rename `.env.example` to `.env` and set the `AUTH_PASSWORD_HASH` variable to the hash you generated in the previous step.
5. Run `node app.js`
6. Open `http://localhost:3000` in your browser


## Stacks
- Node.js
- Alpine.js
- Tailwind CSS
- Express.js