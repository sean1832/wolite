const wol = require("wol");

/**
 * Sends a Wake-on-LAN magic packet.
 * @param {string} mac - The MAC address of the target device.
 * @param {function} callback - Callback function to handle the result.
 */
function sendMagicPacket(mac, callback) {
  wol.wake(mac, function (error) {
    if (error) {
      console.error("Error sending magic packet:", error);
      return callback(error);
    }
    callback(null, `Magic packet sent to ${mac}`);
  });
}

module.exports = sendMagicPacket;
