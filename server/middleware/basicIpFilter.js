const GetErrorPage = require("../utils/serveErrorPage");
const LogConsole = require("../utils/logging");
const config = require("../../data/config.json");

function basicIpAuth(req, res, next) {
  // Get the client IP address. If your server is behind a proxy, you might need additional configuration.
  let clientIp = req.ip;

  // In some cases, IPv4 addresses may be prefixed with "::ffff:" when coming from an IPv6 stack.
  if (clientIp.startsWith("::ffff:")) {
    clientIp = clientIp.substring(7);
  }

  let allowedOrigins;
  try {
    allowedOrigins = config.ALLOWED_ORIGINS;
  } catch (error) {
    LogConsole(
      "error",
      `Error parsing ALLOWED_ORIGINS: ${error}; ALLOWED_ORIGINS: ${JSON.stringify(config.ALLOWED_ORIGINS)}`,
      clientIp
    );
    return res
      .status(500)
      .send(GetErrorPage(500, "Internal Server Error. Check the logs for more details."));
  }

  // If "ALL" is present, allow all IPs.
  if (allowedOrigins.includes("ALL")) {
    LogConsole("info", "IP allowed (all IPs permitted 'ALL').", clientIp);
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
