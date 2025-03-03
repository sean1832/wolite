#!/usr/bin/env node
const fs = require("fs");
const readline = require("readline");
const bcrypt = require("bcrypt");
const otp = require("otpauth");
const createOTP = require("./server/utils/otp");

// Simple CLI argument parser for --key value or boolean flags.
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

// Utility: Generate a random string of the given length.
function generateRandomChar(length) {
  const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-+_&^%$!~";
  let result = "";
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * characters.length));
  }
  return result;
}

// Write the .env file using the provided config object.
function writeEnvFile(config) {
  const envContent = `# Security configs
AUTH_USER="${config.username}"
AUTH_PASSWORD_HASH="${config.hashedPassword}"
SESSION_SECRET="${config.sessionSecret}"

# (optional) OTP configs
ENABLE_OTP=${config.enableOTP}
OTP_SECRET="${config.otpSecret}"
OTP_URI="${config.otpURI}" # OTP URI for QR code

# Server configs
PORT=${config.port}
ENABLE_LOG=${config.enableLog}
ENABLE_IP_FILTER=true
ALLOWED_ORIGINS=${JSON.stringify(config.allowedOrigins)}
COOKIE_LIFETIME=604800000
`;
  fs.writeFileSync(".env", envContent);
  console.log("\n.env file generated successfully.");
}

// Interactive mode using readline.
async function interactiveMode() {
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: true,
  });
  const askQuestion = (query) =>
    new Promise((resolve) => rl.question(query, (answer) => resolve(answer)));

  try {
    let username = await askQuestion("Enter your username: ");
    if (!username) {
      console.error("Username cannot be empty.");
      rl.close();
      process.exit(1);
    }

    let password = await askQuestion("Enter your password: ");
    if (!password) {
      console.error("Password cannot be empty.");
      rl.close();
      process.exit(1);
    }

    let enableOTP = await askQuestion("Enable OTP? (yes/[no]): ");
    enableOTP = enableOTP.trim().toLowerCase() === "yes" ? true : false;

    let port = await askQuestion("Enter the port number (3000): ");
    port = port || "3000";

    let allowedOrigins = await askQuestion(
      "Enter the allowed origins (comma separated, default: ::1,127.0.0.1): "
    );
    if (allowedOrigins.trim() === "") {
      allowedOrigins = ["::1", "127.0.0.1"];
    } else {
      allowedOrigins = allowedOrigins.split(",").map((origin) => origin.trim());
      // Ensure defaults are present.
      if (!allowedOrigins.includes("::1")) allowedOrigins.push("::1");
      if (!allowedOrigins.includes("127.0.0.1")) allowedOrigins.push("127.0.0.1");
      // Remove empty entries.
      allowedOrigins = allowedOrigins.filter((origin) => origin.length > 0);
    }

    let enableLog = await askQuestion("Enable logging? (yes/[no]): ");
    enableLog = enableLog.trim().toLowerCase() === "yes" ? true : false;

    const hashedPassword = await bcrypt.hash(password, 10);
    const sessionSecret = generateRandomChar(32);

    let otpSecret, otpURI;
    if (enableOTP) {
      otpSecret = new otp.Secret({ size: 20 }).base32;
      const totp = createOTP(otpSecret);
      otpURI = totp.toString();
    } else {
      otpSecret = "N/A";
      otpURI = "";
    }

    // Ask user if they want to write the .env file.
    let writeEnv = await askQuestion("\nSave to .env file? ([yes]/no): ");
    if (writeEnv.trim().toLowerCase() === "no") {
      console.log("Not writing .env file. Exiting.");
      rl.close();
      process.exit(0);
    }

    const config = {
      username,
      hashedPassword,
      sessionSecret,
      enableOTP,
      otpSecret,
      otpURI,
      port,
      allowedOrigins,
      enableLog,
    };

    writeEnvFile(config);

    console.log("\nEnvironment variables have been set.");
    console.log("You can now run the server with the following command:");
    console.log("\nnpm run start\n");
    rl.close();
    process.exit(0);
  } catch (error) {
    console.error("An error occurred:", error);
    rl.close();
    process.exit(1);
  }
}

// CLI mode (non-interactive) that uses command-line parameters.
async function cliMode(options) {
  // Show usage help if requested.
  if (options.help || options.h) {
    console.log(`Usage: setup_cli.js --username <username> --password <password> [options]
Options:
  --username         Required. Username for authentication.
  --password         Required. Plain text password (will be hashed).
  --enable-otp       Optional. Include this flag to enable OTP.
  --port             Optional. Server port number (default: 3000).
  --allowed-origins  Optional. Comma separated list of allowed origins (default: "::1,127.0.0.1").
  --enable-log       Optional. Include this flag to enable logging.
  --help, -h         Show this help message.
`);
    process.exit(0);
  }

  // Validate required parameters.
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

  const enableOTP = options["enable-otp"] || false;
  const port = options.port || "3000";
  let allowedOrigins = options["allowed-origins"]
    ? options["allowed-origins"].split(",").map((origin) => origin.trim())
    : [];
  allowedOrigins.push("::1");
  allowedOrigins.push("127.0.0.1");
  allowedOrigins = allowedOrigins.filter((origin) => origin.length > 0);
  const enableLog = options["enable-log"] || false;

  const hashedPassword = await bcrypt.hash(password, 10);
  const sessionSecret = generateRandomChar(32);

  let otpSecret, otpURI;
  if (enableOTP) {
    otpSecret = new otp.Secret({ size: 20 }).base32;
    const totp = createOTP(otpSecret);
    otpURI = totp.toString();
  } else {
    otpSecret = "N/A";
    otpURI = "";
  }

  const config = {
    username,
    hashedPassword,
    sessionSecret,
    enableOTP,
    otpSecret,
    otpURI,
    port,
    allowedOrigins,
    enableLog,
  };

  writeEnvFile(config);

  console.log("\nEnvironment variables have been set.");
  console.log("You can now run the server with the following command:");
  console.log("\nnpm run start\n");
  process.exit(0);
}

// Main: choose interactive or CLI mode based on arguments.
async function main() {
  const options = parseArgs();
  if (process.argv.slice(2).length === 0) {
    // No arguments passed, use interactive Q&A mode.
    await interactiveMode();
  } else {
    await cliMode(options);
  }
}

main().catch((err) => {
  console.error("An error occurred:", err);
  process.exit(1);
});
