const wol = require("wol");
const LogConsole = require("../utils/logging");

/**
 * Sends a Wake-on-LAN magic packet.
 * @param {string} mac - The MAC address of the target device.
 * @param {function} callback - Callback function to handle the result.
 */
function sendMagicPacket(mac, callback) {
  wol.wake(mac, function (error) {
    if (error) {
      LogConsole("error", `Error sending magic packet: ${error}`);
      return callback("Error sending magic packet. Check the logs for more details.");
    }
    callback(null, `Magic packet sent.`);
  });
}

module.exports = sendMagicPacket;
