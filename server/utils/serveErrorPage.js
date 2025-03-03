const fs = require("fs");
const path = require("path");

function GetErrorPage(code, message) {
  const errorFilePath = path.join(__dirname, "..", "..", "public", "error.html");
  let errorHtml = fs.readFileSync(errorFilePath, "utf8");
  errorHtml = errorHtml.replace("{{errorCode}}", code);
  errorHtml = errorHtml.replace("{{errorMessage}}", message);

  return errorHtml;
}

module.exports = GetErrorPage;
