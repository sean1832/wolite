const readline = require("readline");
const bcrypt = require("bcrypt");
const otp = require("otpauth");
const getOTP = require("./server/utils/otp");

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  terminal: true,
});

function askQuestion(query) {
  return new Promise((resolve) => {
    rl.question(query, (answer) => resolve(answer));
  });
}

async function hashPassword(password) {
  try {
    return await bcrypt.hash(password, 10);
  } catch (error) {
    console.error("Error hashing password:", error);
    return null;
  }
}

async function generateRandomChar(length) {
  const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-+_&^%$!~";
  let result = "";
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * characters.length));
  }
  return result;
}

/**
 * Main function to collect user details.
 */
async function main() {
  try {
    const username = await askQuestion("Enter your username: ");
    if (!username) {
      console.error("Username cannot be empty.");
      rl.close();
    }
    const password = await askQuestion("Enter your password: ");
    if (!password) {
      console.error("Password cannot be empty.");
      rl.close();
    }

    let enableOTP = await askQuestion("Enable OTP? (yes/[no]): ");
    enableOTP = enableOTP.toLowerCase() === "yes" ? "true" : "false";

    let port = await askQuestion("Enter the port number (3000): ");
    port = port || "3000";

    let allowedOrigins = await askQuestion("Enter the allowed origins (::1): ");
    allowedOrigins = allowedOrigins.split(",").map((origin) => origin.trim());
    allowedOrigins.push("::1");
    allowedOrigins.push("127.0.0.1");
    // make sure no empty string
    allowedOrigins = allowedOrigins.filter((origin) => origin.length > 0);

    let enableLog = await askQuestion("Enable logging? (yes/[no]): ");
    enableLog = enableLog.toLowerCase() === "yes" ? "true" : "false";

    const hashedPassword = await hashPassword(password);
    const sessionSecret = await generateRandomChar(32);

    let otpSecret;
    let otpURI;
    if (enableOTP === "true") {
      otpSecret = new otp.Secret({ size: 20 }).base32;
      const otpInstance = getOTP(otpSecret);
      otpURI = otpInstance.toString();
    } else {
      otpSecret = "N/A";
    }

    console.log("\nEnvironment Variables:");
    console.log(`AUTH_USER="${username}"`);
    console.log(`AUTH_PASSWORD_HASH="${hashedPassword}"`);
    console.log(`SESSION_SECRET="${sessionSecret}"`);
    console.log(`ENABLE_OTP=${enableOTP}`);
    console.log(`OTP_URI="${otpURI}"`);
    console.log(`OTP_SECRET="${otpSecret}"`);

    console.log(`PORT=${port}`);
    console.log(`ENABLE_LOG=${enableLog}`);
    console.log(`ENABLE_IP_FILTER=true`);
    console.log(`ALLOWED_ORIGINS=${JSON.stringify(allowedOrigins)}`);
    console.log(`COOKIE_LIFETIME=604800000`);

    const writeEnvFile = await askQuestion("\nSave to .env file? ([yes]/no): ");
    if (writeEnvFile.toLowerCase() != "no") {
      console.log("\nWriting to .env file...");
      const fs = require("fs");

      fs.writeFileSync(".env", `# Security configs\n`);
      fs.appendFileSync(".env", `AUTH_USER="${username}"\n`);
      fs.appendFileSync(".env", `AUTH_PASSWORD_HASH="${hashedPassword}"\n`);
      fs.appendFileSync(".env", `SESSION_SECRET="${sessionSecret}"\n\n`);

      fs.appendFileSync(".env", `# (optional) OTP configs\n`);
      fs.appendFileSync(".env", `ENABLE_OTP=${enableOTP}\n`);
      fs.appendFileSync(".env", `OTP_SECRET="${otpSecret}"\n`);
      fs.appendFileSync(".env", `# (Optional) You can add this URI to your OTP app\n`);
      fs.appendFileSync(".env", `OTP_URI="${otpURI}"\n\n`);

      fs.appendFileSync(".env", `# Server configs\n`);
      fs.appendFileSync(".env", `PORT=${port}\n`);
      fs.appendFileSync(".env", `ENABLE_LOG=${enableLog}\n`);
      fs.appendFileSync(".env", `ENABLE_IP_FILTER=true # disable to allow all IPs\n`);
      fs.appendFileSync(
        ".env",
        `ALLOWED_ORIGINS=${JSON.stringify(
          allowedOrigins
        )} # add more as needed in JSON array format\n`
      );
      fs.appendFileSync(".env", `COOKIE_LIFETIME=604800000 # 7 days in milliseconds\n`);
    }

    console.log("\nEnvironment variables have been set.");
    console.log("You can now run the server with the following command:");
    console.log("\nnpm run start\n");
    rl.close();
  } catch (error) {
    console.error("An error occurred:", error);
    rl.close();
  }
}

main();
