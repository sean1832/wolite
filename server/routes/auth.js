const express = require("express");
const path = require("path");
const bcrypt = require("bcrypt");

const router = express.Router();

// GET /auth - serve the login page
router.get("/", (req, res) => {
  res.sendFile(path.join(__dirname, "..", "..", "public", "login.html"));
});

// POST /auth - process the login form submission
router.post("/", (req, res) => {
  const { username, password } = req.body;

  // Check username against environment variable
  if (username === process.env.AUTH_USER) {
    bcrypt.compare(password, process.env.AUTH_PASSWORD_HASH, (err, result) => {
      if (result) {
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
