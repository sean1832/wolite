require("dotenv").config(); // Load environment variables
const fs = require("fs");
const log_dir = __dirname + "/../../logs";
if (!fs.existsSync(log_dir)) {
  fs.mkdirSync(log_dir);
}
const date = new Date();
const server_start_time =
  `${date.getUTCFullYear()}` +
  `${date.getUTCMonth()}` +
  `${date.getUTCDate()}` +
  `${date.getUTCHours()}` +
  `${date.getUTCMinutes()}` +
  `${date.getUTCSeconds()}`;
const log_file = fs.createWriteStream(log_dir + `/${server_start_time}.log`, { flags: "w" });

const LogToFile = (ip, content) => {
  if (process.env.ENABLE_LOG == "true") {
    const current_time_utc = new Date().toUTCString();
    log_file.write(`[${current_time_utc}] [${ip}] ${content}\n`);
  }
};

module.exports = LogToFile;
