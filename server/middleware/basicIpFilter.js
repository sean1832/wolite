const GetErrorPage = require("../utils/serveErrorPage");
const LogConsole = require("../utils/logging");
require("dotenv").config();

function basicIpAuth(req, res, next) {
  // Get the client IP address. If your server is behind a proxy, you might need additional configuration.
  let clientIp = req.ip;

  // In some cases, IPv4 addresses may be prefixed with "::ffff:" when coming from an IPv6 stack.
  if (clientIp.startsWith("::ffff:")) {
    clientIp = clientIp.substring(7);
  }

  let allowedOrigins;
  try {
    allowedOrigins = JSON.parse(process.env.ALLOWED_ORIGINS);
  } catch (error) {
    LogConsole(
      "error",
      `Error parsing ALLOWED_ORIGINS: ${error}; ALLOWED_ORIGINS: ${process.env.ALLOWED_ORIGINS}`,
      clientIp
    );
    return res
      .status(500)
      .send(GetErrorPage(500, "Internal Server Error. Check the logs for more details."));
  }

  // If "ANY" is present, allow all IPs.
  if (allowedOrigins.includes("ANY")) {
    LogConsole("info", "IP allowed (all IPs permitted 'ANY').", clientIp);
    return next();
  }

  // Check if the client's IP is in the allowedIPs list
  if (!allowedOrigins.includes(clientIp)) {
    LogConsole("warn", "Access denied. IP not allowed.", clientIp);
    return res.status(403).send(GetErrorPage(403, "Access denied. IP not allowed."));
  }
  LogConsole("info", "IP allowed.", clientIp);
  next();
}

module.exports = basicIpAuth;
