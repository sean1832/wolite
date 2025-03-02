const express = require("express");
const path = require("path");
const bcrypt = require("bcrypt");
require("dotenv").config();

const router = express.Router();

// GET /auth/config - serve the login page
router.get("/", (req, res) => {
  res.json({ enableOtp: process.env.ENABLE_OTP === "true" });
});

module.exports = router;
