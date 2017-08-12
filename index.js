const Interpreter = require("./lib/interpreter"),
  Scopes = require("./lib/scope"),
  fs = require("fs"),
  dashdash = require("dashdash");

function die(message) {
  console.warn(message);
  process.exit(1);
}

var optParser = dashdash.createParser({
  options: [
    {
      names:["help", "h"],
      type: "bool",
      help: "Print usage and exit"
    },
    {
      names:["output", "o"],
      type: "string",
      help: "Output file",
      helpArg: "FILE",
      default: "entities.json"
    },
    {
      names: ["check", "c"],
      type: "bool",
      help: "Check syntax without generating output"
    }
  ]
});


function usage(parser) {
  return "USAGE: panther [OPTIONS] input-file\n\nOPTIONS:\n" + parser.help({includeEnv: true}).trimRight();
}

var opts;

try {
  opts = optParser.parse(process.argv);
} catch(e) {
  console.warn("FATAL: " + e.message + "\n");
  die(usage(optParser));
}

if (opts.help) {
  console.log(usage(optParser));
  process.exit(0);
}

if (opts._args.length !== 1) {
  console.warn("FATAL: Requires an input file\n");
  die(usage(optParser));
}

var inputFile = opts._args[0];
var outputFile = opts.output;

var i = new Interpreter();

if (opts.check) {
  try {
    i.loadFile(inputFile, Scopes.newRootScope(opts.check));
    console.log("Syntax OK");
    process.exit(0);
  } catch(e) {
    die(`Syntax FAILED:\n${e.message}`);
  }
} else {
  i.loadFile(inputFile, Scopes.newRootScope());
  fs.writeFileSync(outputFile, JSON.stringify(i.output, null, 2));
}
