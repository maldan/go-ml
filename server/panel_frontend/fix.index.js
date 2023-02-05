const fs = require("fs");

let file = fs.readFileSync("./dist/index.html", "utf8");
file = file.replace("/assets/", "assets/");
file = file.replace("/assets/", "assets/");
fs.writeFileSync("./dist/index.html", file);
