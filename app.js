require("dotenv").config();
const express = require("express");
const path = require("path");
const session = require("express-session");
const bodyParser = require("body-parser");

// Import the routes
const wakeRouter = require("./server/routes/wake");
const authRouter = require("./server/routes/auth");
const sessionRouter = require("./server/routes/session");
const authConfigRouter = require("./server/routes/authConfig");
const GetDynamicPage = require("./server/utils/serveDynamicPage");

const app = express();
const port = process.env.PORT || 3000;

// Parse URL-encoded bodies (for form submissions)
app.use(bodyParser.urlencoded({ extended: true }));

// Configure session management
app.use(
  session({
    secret: process.env.SESSION_SECRET,
    resave: false,
    saveUninitialized: false,
    cookie: {
      maxAge: Number(process.env.COOKIE_LIFETIME),
    },
  })
);

// Serve static assets (CSS, images, etc.)
app.use(express.static(path.join(__dirname, "public"), { index: false }));

// Main route: check authentication before serving the app
app.get("/", (req, res) => {
  if (!req.session.isAuthenticated) {
    return res.redirect("/auth");
  }

  versionData = {
    version: process.env.npm_package_version,
  };

  // Serve the main application page if authenticated
  res.send(GetDynamicPage("index.html", versionData));
});

// Mount the session status route at /api/session-status
app.use("/api/session-status", sessionRouter);

// Mount the authentication routes at /auth
app.use("/auth", authRouter);

app.use("/auth/config", authConfigRouter);

// Mount the wake route (protected) at /api/wake
app.use("/api/wake", wakeRouter);

app.listen(port, () => {
  console.log(`Server listening at http://localhost:${port}`);
});
