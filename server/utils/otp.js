const otp = require("otpauth");

function createOTP(secret) {
  let totp = new otp.TOTP({
    issuer: "WOLITE",
    label: "WOLITE",
    algorithm: "SHA256",
    digits: 6,
    period: 30,
    secret: secret,
  });
  return totp;
}

module.exports = createOTP;
