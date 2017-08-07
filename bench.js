const hyperid = require("hyperid");
const shortid = require("shortid");
var i;

const hyper = hyperid();

const iters = 1000000;

console.time(`${iters.toLocaleString()} shortid`)
for (i = 0; i < iters; ++i) {
  shortid.generate();
}
console.timeEnd(`${iters.toLocaleString()} shortid`)

console.time(`${iters.toLocaleString()} hyperid`)
for (i = 0; i < iters; ++i) {
  hyper();
}
console.timeEnd(`${iters.toLocaleString()} hyperid`)

