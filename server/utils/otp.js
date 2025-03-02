const otp = require("otpauth");

function getOTP(secret) {
  let totp = new otp.TOTP({
    issuer: "WOLITE",
    label: "WOLITE:SECRET",
    algorithm: "SHA256",
    digits: 6,
    period: 30,
    secret: secret,
  });
  return totp;
}

module.exports = getOTP;
