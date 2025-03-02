// server/routes/session.js
const express = require("express");
const router = express.Router();

router.get("/", (req, res) => {
  if (req.session && req.session.isAuthenticated) {
    return res.json({ authenticated: true });
  }
  res.json({ authenticated: false });
});

module.exports = router;
