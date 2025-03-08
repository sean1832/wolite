const express = require("express");
const config = require("../../data/config.json");

const router = express.Router();

// GET /auth/config - serve the login page
router.get("/", (req, res) => {
  res.json({ enableOtp: config.ENABLE_OTP === true });
});

module.exports = router;
