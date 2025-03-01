const readline = require("readline");
const Writable = require("stream").Writable;
const bcrypt = require("bcrypt");

var mutableStdout = new Writable({
  write: function (chunk, encoding, callback) {
    if (!this.muted) process.stdout.write(chunk, encoding);
    callback();
  },
});

mutableStdout.muted = false;

var rl = readline.createInterface({
  input: process.stdin,
  output: mutableStdout,
  terminal: true,
});

rl.question("Password: ", function (password) {
  bcrypt.hash(password, 10, function (err, hash) {
    if (err) {
      console.error("Error hashing password:", err);
      return rl.close();
    }
    console.log("\n");
    console.log(hash);
    rl.close();
  });
  rl.close();
});

mutableStdout.muted = true;
