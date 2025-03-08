#!/usr/bin/env node
// This script is used to generate `config.json` based on the environment variables provided.
// It is intended to be run during the Docker container setup process.
// If you are running the server outside of a Docker container, you can use the CLI setup script instead.
require("dotenv").config();
const fs = require("fs");
const bcrypt = require("bcrypt");
const otp = require("otpauth");
const createOTP = require("./server/utils/otp");

function writeConfig(path, data) {
  fs.writeFileSync(path, JSON.stringify(data, null, 4));
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

// default configuration
let config = {
  AUTH_USER: process.env.USERNAME,
  AUTH_PASSWORD_HASH: bcrypt.hash(process.env.PASSWORD, 10),
  SESSION_SECRET: generateRandomChar(32),
  SESSION_LIFETIME: 1800000,
  ENABLE_OTP: process.env.ENABLE_OTP === "true",
  OTP_SECRET: "N/A",
  OTP_URI: "",
  PORT: Number(process.env.PORT) || 3000,
  ALLOWED_ORIGINS: ["::1", "127.0.0.1"],
  COOKIE_LIFETIME: 604800000,
};

async function main() {
  // OTP configuration
  if (config.ENABLE_OTP) {
    let otpSecret;
    let totp;
    if (process.env.OTP_URI) {
      // If the OTP URI is provided, parse it and use the secret
      totp = otp.URI.parse(process.env.OTP_URI);
      otpSecret = totp.secret.base32;
    } else if (process.env.OTP_SECRET) {
      // If the OTP secret is provided, use it
      otpSecret = process.env.OTP_SECRET;
      if (!process.env.OTP_ALGORITHM || !process.env.OTP_DIGITS || !process.env.OTP_PERIOD) {
        throw new Error(
          "Missing required OTP configuration. Required: OTP_ALGORITHM, OTP_DIGITS, OTP_PERIOD if OTP_SECRET is provided."
        );
      }
      totp = createOTP(
        otpSecret,
        process.env.OTP_ALGORITHM,
        process.env.OTP_DIGITS,
        process.env.OTP_PERIOD
      );
    } else {
      // Otherwise, generate a new secret
      otpSecret = new otp.Secret({ size: 20 }).base32;
      totp = createOTP(otpSecret);
    }
    const otpURI = totp.toString();

    config.OTP_SECRET = otpSecret;
    config.OTP_URI = otpURI;
  }

  // Allowed origins
  if (process.env.ALLOWED_ORIGINS) {
    const envAllowedOrigins = process.env.ALLOWED_ORIGINS.split(",");
    const defaults = config.ALLOWED_ORIGINS;
    envAllowedOrigins.forEach((origin) => {
      const trimmedOrigin = origin.trim();
      if (!defaults.includes(trimmedOrigin)) {
        config.ALLOWED_ORIGINS.push(trimmedOrigin);
      }
    });
  }

  writeConfig("data/config.json", config);
}

main().catch((err) => {
  console.error("An error occurred:", err);
  process.exit(1);
});
