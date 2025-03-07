#!/usr/bin/env node
// This script is used to generate `config.json` based on user input.
// Run `npm run setup` to start the interactive setup process.
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
  const characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&()*+,-./:;<=>?@[]^_`{|}~";
  let result = "";
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * characters.length));
  }
  return result;
}

/**
 * Process allowed origins into an array. Defaults to [ "::1", "127.0.0.1" ].
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
 */
async function buildConfig(input) {
  const hashedPassword = await bcrypt.hash(input.password, 10);
  const sessionSecret = generateRandomChar(32);

  let otpSecret, otpURI;
  if (input.enableOTP) {
    // If OTP URI is provided, use it
    if (input.otpURI && input.otpURI.trim() !== "") {
      const totp = otp.URI.parse(input.otpURI);
      otpSecret = totp.secret.base32;
      otpURI = totp.toString();
    } else {
      // Otherwise, generate a new secret and OTP
      otpSecret = new otp.Secret({ size: 20 }).base32;
      const totp = createOTP(otpSecret);
      otpURI = totp.toString();
    }
  } else {
    otpSecret = "N/A";
    otpURI = "";
  }

  return {
    AUTH_USER: input.username,
    AUTH_PASSWORD_HASH: hashedPassword,
    SESSION_SECRET: sessionSecret,
    SESSION_LIFETIME: 1800000, // 30 minutes in ms
    ENABLE_OTP: input.enableOTP,
    OTP_SECRET: otpSecret,
    OTP_URI: otpURI,
    PORT: Number(input.port) || 3000,
    ALLOWED_ORIGINS: processAllowedOrigins(input.allowedOrigins),
    COOKIE_LIFETIME: 604800000, // 7 days in ms
  };
}

/**
 * Write the configuration object to config.json.
 */
function writeConfig(config) {
  fs.writeFileSync("data/config.json", JSON.stringify(config, null, 4));
  console.log("\nconfig.json generated successfully.");
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

    // Only ask for OTP URI if OTP is enabled
    let otpURIInput = "";
    if (enableOTP) {
      otpURIInput = await askQuestion("Enter OTP URI (leave blank to generate new one): ");
    }

    const port = (await askQuestion("Enter the port number (3000): ")) || "3000";

    const originsAnswer = await askQuestion(
      "Enter the allowed origins (comma separated, default: ::1,127.0.0.1): "
    );

    const configInput = {
      username,
      password,
      enableOTP,
      otpURI: otpURIInput,
      port,
      allowedOrigins: originsAnswer,
    };

    const config = await buildConfig(configInput);

    const writeConfigAnswer = await askQuestion("\nSave to config.json? ([yes]/no): ");
    if (writeConfigAnswer.trim().toLowerCase() === "no") {
      console.log("Not writing config.json. Exiting.");
      rl.close();
      process.exit(0);
    }

    writeConfig(config);

    console.log("\nConfiguration has been set.");
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
  --otp-uri          Optional. Existing OTP URI to use (if OTP is enabled).
  --port             Optional. Server port number (default: 3000).
  --allowed-origins  Optional. Comma separated list of allowed origins (default: "::1,127.0.0.1").
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
    otpURI: options["otp-uri"] || "",
    port: options.port || "3000",
    allowedOrigins: options["allowed-origins"] || "",
  };

  const config = await buildConfig(configInput);
  writeConfig(config);

  console.log("\nConfiguration has been set.");
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
