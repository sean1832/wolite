require("dotenv").config(); // Load environment variables
const express = require("express");
const path = require("path");

// Import the wake endpoint router
const wakeRouter = require("./server/routes/wake");

const app = express();
const port = process.env.PORT || 3000;

// Serve static files (e.g., index.html) from the public directory
app.use(express.static(path.join(__dirname, "public")));

// Mount the wake router at the /wake path
app.use("/wake", wakeRouter);

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
