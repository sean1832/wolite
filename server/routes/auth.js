const express = require("express");
const path = require("path");
const bcrypt = require("bcrypt");
require("dotenv").config();
const createOTP = require("../utils/otp");

const router = express.Router();

// GET /auth - serve the login page
router.get("/", (req, res) => {
  res.sendFile(path.join(__dirname, "..", "..", "public", "login.html"));
});

// POST /auth - process the login form submission
router.post("/", (req, res) => {
  const { username, password, totp } = req.body;

  // Check username against environment variable
  if (username === process.env.AUTH_USER) {
    bcrypt.compare(password, process.env.AUTH_PASSWORD_HASH, (err, passwordMatch) => {
      if (passwordMatch) {
        // Check OTP if enabled
        if (process.env.ENABLE_OTP === "true") {
          try {
            if (!totp) {
              return res.redirect("/auth?error=OTP%20required.");
            }
            if (process.env.OTP_SECRET === "N/A" || !process.env.OTP_SECRET) {
              return res.redirect("/auth?error=OTP%20is%20enabled%20but%20not%20configured.");
            }
            const totpInstance = createOTP(process.env.OTP_SECRET);
            const delta = totpInstance.validate({ token: totp, window: 1 });
            if (delta === null) {
              return res.redirect("/auth?error=Invalid%20OTP.");
            }
          } catch (error) {
            console.error("TOTP validation error:", error);
            return res.redirect("/auth?error=Error%20validating%20TOTP.%20(Internal%20Error)");
          }
        }

        // Mark the session as authenticated
        req.session.isAuthenticated = true;
        // Redirect to the main page (index.html)
        return res.redirect("/");
      } else {
        return res.redirect("/auth?error=Invalid%20username%20or%20password.");
      }
    });
  } else {
    return res.redirect("/auth?error=Invalid%20username%20or%20password.");
  }
});

module.exports = router;
