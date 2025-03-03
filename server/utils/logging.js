/**
 * Logs a message with the specified level.
 *
 * @param {"log"|"warn"|"error"|"info"|"debug"} [level="log"] - The log level.
 * @param {string} content - The content to log.
 * @param {string} ip - The IP address.
 */
const LogConsole = (level = "info", content, ip) => {
  if (!content) {
    throw new Error("Content are required.");
  }
  let msg;
  if (!ip) {
    msg = `${content}\n`;
  } else {
    msg = `[${ip}] ${content}\n`;
  }

  if (level === "error") {
    console.error(msg);
    return;
  } else if (level === "warn") {
    console.warn(msg);
    return;
  } else if (level === "info") {
    console.info(msg);
    return;
  } else if (level === "debug") {
    console.debug(msg);
    return;
  } else if (level === "log") {
    console.log(msg);
  } else {
    throw new Error("Invalid log level.");
  }
};

module.exports = LogConsole;
