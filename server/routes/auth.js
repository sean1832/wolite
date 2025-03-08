const express = require("express");
const path = require("path");
const bcrypt = require("bcrypt");
require("dotenv").config();
const createOTP = require("../utils/otp");
const LogConsole = require("../utils/logging");
const basicIpFilter = require("../middleware/basicIpFilter");
const config = require("../../data/config.json");

const router = express.Router();

// GET /auth - serve the login page
router.get("/", (req, res) => {
  res.sendFile(path.join(__dirname, "..", "..", "public", "login.html"));
});

router.use(basicIpFilter);

// POST /auth - process the login form submission
router.post("/", (req, res) => {
  const { username, password, totp } = req.body;

  // Check username against environment variable
  if (username === config.AUTH_USER) {
    bcrypt.compare(password, config.AUTH_PASSWORD_HASH, (err, passwordMatch) => {
      if (passwordMatch) {
        // Check OTP if enabled
        if (config.ENABLE_OTP === true) {
          try {
            if (!totp) {
              LogConsole("warn", "OTP required.", req.ip);
              return res.redirect("/auth?error=OTP%20required.");
            }
            if (config.OTP_SECRET === "N/A" || !config.OTP_SECRET) {
              LogConsole("warn", "OTP is enabled but not configured.", req.ip);
              return res.redirect("/auth?error=OTP%20is%20enabled%20but%20not%20configured.");
            }
            const totpInstance = createOTP(config.OTP_SECRET);
            const delta = totpInstance.validate({ token: totp, window: 1 });
            if (delta === null) {
              LogConsole("warn", "Invalid OTP provided.", req.ip);
              return res.redirect("/auth?error=Invalid%20OTP.");
            }
          } catch (error) {
            LogConsole("error", `Error validating TOTP: ${error}`, req.ip);
            return res.redirect("/auth?error=Error%20validating%20TOTP.%20(Internal%20Error)");
          }
        }

        // Mark the session as authenticated
        req.session.isAuthenticated = true;
        LogConsole("info", "User authenticated.", req.ip);
        // Redirect to the main page (index.html)
        return res.redirect("/");
      } else {
        LogConsole("warn", "Invalid password.", req.ip);
        return res.redirect("/auth?error=Invalid%20username%20or%20password.");
      }
    });
  } else {
    LogConsole("warn", "Invalid username.", req.ip);
    return res.redirect("/auth?error=Invalid%20username%20or%20password.");
  }
});

module.exports = router;
