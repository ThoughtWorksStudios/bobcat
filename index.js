const fs = require("fs");
const Parser = require("./dsl");

var path = "person.lang";

if (process.argv.length > 2) {
  path = process.argv[2];
}

var input = fs.readFileSync(path, "utf-8");
console.log("Input: ", input, "\n");

console.log("Parsed: ", JSON.stringify(Parser.parse(input), null, 2));
