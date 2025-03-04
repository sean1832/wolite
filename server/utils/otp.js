const otp = require("otpauth");

function createOTP(secret, algorithm = "SHA256", digits = 6, period = 30) {
  let totp = new otp.TOTP({
    issuer: "WOLITE",
    label: "WOLITE",
    algorithm: algorithm,
    digits: digits,
    period: period,
    secret: secret,
  });
  return totp;
}

module.exports = createOTP;
