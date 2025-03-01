require("dotenv").config(); // Load environment variables
const bcrypt = require("bcrypt");
const LogToFile = require("../utils/logging");

// Custom basic auth middleware
const basicPasswordAuth = (req, res, next) => {
  const authHeader = req.headers.authorization;
  if (!authHeader) {
    res.setHeader("WWW-Authenticate", "Basic");
    return res.status(401).send("Authentication required.");
  }

  const [scheme, base64] = authHeader.split(" ");
  if (scheme !== "Basic") {
    LogToFile(req.ip, "Invalid authentication scheme. Only 'Basic' is supported.");
    return res.status(401).send("Invalid authentication scheme.");
  }

  const credentials = Buffer.from(base64, "base64").toString();
  const [username, password] = credentials.split(":");
  if (!username || !password) {
    return res.status(401).send("No username or password provided.");
  }

  // Check the username from environment variable
  if (username !== process.env.AUTH_USER) {
    LogToFile(req.ip, `Unauthorized. Incorrect username: ${username}`);
    return res.status(401).send("Unauthorized. Incorrect password or username.");
  }

  // Compare provided password with stored hash
  bcrypt.compare(password, process.env.AUTH_PASSWORD_HASH, (err, result) => {
    if (err || !result) {
      LogToFile(
        req.ip,
        `Unauthorized. Incorrect password for user ${username}. Password attempt: ${password}`
      );
      return res.status(401).send("Unauthorized. Incorrect password or username.");
    }
    next();
  });
};

module.exports = basicPasswordAuth;
