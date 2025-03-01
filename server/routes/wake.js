const express = require("express");
const router = express.Router();

const basicPasswordAuth = require("../middleware/basicPasswordAuth");
const basicIpFilter = require("../middleware/basicIpFilter");
const sendMagicPacket = require("../logic/sendMagicPacket");
const LogToFile = require("../utils/logging");

// Conditionally apply IP filter if enabled
if (process.env.ENABLE_IP_FILTER === "true") {
  router.use(basicIpFilter);
}

// Apply basic password authentication middleware
router.use(basicPasswordAuth);

// Define the /wake GET endpoint
router.get("/", (req, res) => {
  const clientIp = req.ip;
  const mac = req.query.mac;

  if (!mac) {
    LogToFile(clientIp, "MAC address not provided.");
    return res.status(400).type("text/plain").send("MAC address not provided.");
  }

  // Validate MAC address format
  const macRegex = /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/;
  if (!macRegex.test(mac)) {
    LogToFile(clientIp, "Invalid MAC address format.");
    return res.status(400).type("text/plain").send("Invalid MAC address format.");
  }

  // Call the magic packet sending function
  sendMagicPacket(mac, (error, result) => {
    if (error) {
      LogToFile(clientIp, `Error sending magic packet: ${error}`);
      return res
        .status(500)
        .type("text/plain")
        .send("Error sending magic packet. Check the logs for more details.");
    }
    LogToFile(clientIp, result);
    res.cookie("mac", mac, { maxAge: 2_592_000_000 }); // 30 days
    res.status(200).type("text/plain").send(result);
  });
});

module.exports = router;
