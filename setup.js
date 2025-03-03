#!/usr/bin/env node

const fs = require("fs");
const readline = require("readline");
const bcrypt = require("bcrypt");
const otp = require("otpauth");
const createOTP = require("./server/utils/otp");

/**
 * Parse CLI arguments in the form --key value or boolean flags.
 */
function parseArgs() {
  const args = process.argv.slice(2);
  const options = {};
  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg.startsWith("--")) {
      const key = arg.slice(2);
      if (i + 1 < args.length && !args[i + 1].startsWith("-")) {
        options[key] = args[i + 1];
        i++;
      } else {
        options[key] = true;
      }
    } else if (arg.startsWith("-")) {
      const key = arg.slice(1);
      if (i + 1 < args.length && !args[i + 1].startsWith("-")) {
        options[key] = args[i + 1];
        i++;
      } else {
        options[key] = true;
      }
    }
  }
  return options;
}

/**
 * Generate a random string of the given length.
 */
function generateRandomChar(length) {
  const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-+_&^%$!~";
  let result = "";
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * characters.length));
  }
  return result;
}

/**
 * Ensure allowed origins are processed into an array with default values.
 */
function processAllowedOrigins(originsInput) {
  if (!originsInput || (typeof originsInput === "string" && originsInput.trim() === "")) {
    return ["::1", "127.0.0.1"];
  }
  let origins = Array.isArray(originsInput)
    ? originsInput
    : originsInput
        .split(",")
        .map((o) => o.trim())
        .filter((o) => o);
  if (!origins.includes("::1")) origins.push("::1");
  if (!origins.includes("127.0.0.1")) origins.push("127.0.0.1");
  return origins;
}

/**
 * Build the configuration object based on provided inputs.
 * It handles password hashing, OTP generation, session secret, etc.
 */
async function buildConfig(input) {
  const hashedPassword = await bcrypt.hash(input.password, 10);
  const sessionSecret = generateRandomChar(32);

  let otpSecret, otpURI;
  if (input.enableOTP) {
    otpSecret = new otp.Secret({ size: 20 }).base32;
    const totp = createOTP(otpSecret);
    otpURI = totp.toString();
  } else {
    otpSecret = "N/A";
    otpURI = "";
  }

  return {
    username: input.username,
    hashedPassword,
    sessionSecret,
    enableOTP: input.enableOTP,
    otpSecret,
    otpURI,
    port: input.port || "3000",
    allowedOrigins: processAllowedOrigins(input.allowedOrigins),
  };
}

/**
 * Write the .env file using the provided configuration.
 */
function writeEnvFile(config) {
  const envContent = `# Security configs
AUTH_USER="${config.username}"
AUTH_PASSWORD_HASH="${config.hashedPassword}"
SESSION_SECRET="${config.sessionSecret}"
SESSION_LIFETIME=1800000 # 30 minutes in milliseconds

# (optional) OTP configs
ENABLE_OTP=${config.enableOTP}
OTP_SECRET="${config.otpSecret}"
OTP_URI="${config.otpURI}" # OTP URI for QR code

# Server configs
PORT=${config.port}
ALLOWED_ORIGINS=${JSON.stringify(config.allowedOrigins)} # add "*" to allow all origins
COOKIE_LIFETIME=604800000 # 7 days in milliseconds
`;
  fs.writeFileSync(".env", envContent);
  console.log("\n.env file generated successfully.");
}

/**
 * Interactive mode: uses readline to prompt the user.
 */
async function interactiveMode() {
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: true,
  });

  const askQuestion = (query) => new Promise((resolve) => rl.question(query, resolve));

  try {
    const username = await askQuestion("Enter your username: ");
    if (!username) {
      console.error("Username cannot be empty.");
      rl.close();
      process.exit(1);
    }

    const password = await askQuestion("Enter your password: ");
    if (!password) {
      console.error("Password cannot be empty.");
      rl.close();
      process.exit(1);
    }

    const otpAnswer = await askQuestion("Enable OTP? (yes/[no]): ");
    const enableOTP = otpAnswer.trim().toLowerCase() === "yes";

    const port = (await askQuestion("Enter the port number (3000): ")) || "3000";

    const originsAnswer = await askQuestion(
      "Enter the allowed origins (comma separated, '*' for any origin, default: ::1,127.0.0.1): "
    );

    const configInput = {
      username,
      password,
      enableOTP,
      port,
      allowedOrigins: originsAnswer,
    };

    const config = await buildConfig(configInput);

    const writeEnvAnswer = await askQuestion("\nSave to .env file? ([yes]/no): ");
    if (writeEnvAnswer.trim().toLowerCase() === "no") {
      console.log("Not writing .env file. Exiting.");
      rl.close();
      process.exit(0);
    }

    writeEnvFile(config);

    console.log("\nEnvironment variables have been set.");
    console.log("You can now run the server with:\n\nnpm run start\n");
    rl.close();
    process.exit(0);
  } catch (error) {
    console.error("An error occurred:", error);
    rl.close();
    process.exit(1);
  }
}

/**
 * CLI mode: uses command-line arguments.
 */
async function cliMode(options) {
  if (options.help || options.h) {
    console.log(`Usage: setup_cli.js --username <username> --password <password> [options]
Options:
  --username         Required. Username for authentication.
  --password         Required. Plain text password (will be hashed).
  --enable-otp       Optional. Include this flag to enable OTP.
  --port             Optional. Server port number (default: 3000).
  --allowed-origins  Optional. Comma separated list of allowed origins. '*' for any origin. (default: "::1,127.0.0.1").
  --help, -h         Show this help message.
`);
    process.exit(0);
  }

  const username = options.username || options.u;
  const password = options.password || options.p;
  if (!username) {
    console.error("Error: --username is required.");
    process.exit(1);
  }
  if (!password) {
    console.error("Error: --password is required.");
    process.exit(1);
  }

  const configInput = {
    username,
    password,
    enableOTP: options["enable-otp"] || false,
    port: options.port || "3000",
    allowedOrigins: options["allowed-origins"] || "",
  };

  const config = await buildConfig(configInput);
  writeEnvFile(config);

  console.log("\nEnvironment variables have been set.");
  console.log("You can now run the server with:\n\nnpm run start\n");
  process.exit(0);
}

/**
 * Main entry point: decide between interactive and CLI mode.
 */
async function main() {
  const options = parseArgs();
  if (process.argv.slice(2).length === 0) {
    await interactiveMode();
  } else {
    await cliMode(options);
  }
}

main().catch((err) => {
  console.error("An error occurred:", err);
  process.exit(1);
});
