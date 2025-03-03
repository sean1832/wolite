const fs = require("fs");
const path = require("path");

function GetDynamicPage(pageName, data) {
  const errorFilePath = path.join(__dirname, "..", "..", "public", pageName);
  let html = fs.readFileSync(errorFilePath, "utf8");
  for (const key in data) {
    html = html.replace(`{{${key}}}`, data[key]);
  }
  return html;
}

module.exports = GetDynamicPage;
