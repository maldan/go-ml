const fs = require("fs");

let file = fs.readFileSync("./dist/index.html", "utf8");
file = file.replace("/assets/index.js", "/debug/panel/js");
file = file.replace("/assets/index.css", "/debug/panel/css");
fs.writeFileSync("./dist/index.html", file);
