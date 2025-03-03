const express = require("express");
const router = express.Router();
const basicIpFilter = require("../middleware/basicIpFilter");
const sendMagicPacket = require("../logic/sendMagicPacket");
const LogConsole = require("../utils/logging");

// Custom middleware to check if the user is authenticated
const customAuth = (req, res, next) => {
  if (req.session.isAuthenticated) {
    return next();
  }
  res.redirect("/auth");
};

// Optionally apply IP filter if enabled
if (process.env.ENABLE_IP_FILTER === "true") {
  router.use(basicIpFilter);
}

// Apply the session authentication middleware
router.use(customAuth);

// GET /wake endpoint
router.get("/", (req, res) => {
  const clientIp = req.ip;
  const mac = req.query.mac;

  if (!mac) {
    LogConsole("warn", "MAC address not provided.", clientIp);
    return res.status(400).type("text/plain").send("MAC address not provided.");
  }

  // Validate MAC address format
  const macRegex = /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/;
  if (!macRegex.test(mac)) {
    LogConsole("warn", "Invalid MAC address format.", clientIp);
    return res.status(400).type("text/plain").send("Invalid MAC address format.");
  }

  // Send the magic packet
  sendMagicPacket(mac, (error, result) => {
    if (error) {
      LogConsole("error", `Error sending magic packet: ${error}`, clientIp);
      return res
        .status(500)
        .type("text/plain")
        .send("Error sending magic packet. Check the logs for more details.");
    }
    LogConsole("info", result, clientIp);
    res.cookie("mac", mac, { maxAge: Number(process.env.COOKIE_LIFETIME) });
    res.status(200).type("text/plain").send(result);
  });
});

module.exports = router;
