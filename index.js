const Interpreter = require("./interpreter");
const scp = require("./scope");
const fs = require("fs");

var path = "examples/example.lang";

if (process.argv.length > 2) {
  path = process.argv[2];
}

var i = new Interpreter();

i.loadFile(path, scp.newRootScope());

fs.writeFileSync("entities.json", JSON.stringify(i.output, null, 2));
