;(function() {
  "use strict";

  const fs = require("fs");

  function Scope(parent) {
    this.symbols = {};
    this.imports = new FileHash();

    this.set = function set(symbol, value) {
      this.symbols[symbol] = value;
    };

    this.resolve = function resolveSymbol(symbol) {
      if ("undefined" !== typeof this.symbols[symbol]) {
        return this.symbols[symbol];
      }

      if (parent) {
        return parent.resolve(symbol);
      }

      throw new Error(`Failed to resolve symbol ${JSON.stringify(symbol)}`);
    };
  }

  function FileHash() {
    var files = {};
    this.hasSeen = function hasSeen(fspath) {
      return !!files[fs.realpathSync(fspath)];
    };

    this.markSeen = function markSeen(fspath) {
      files[fs.realpathSync(fspath)] = true;
    };
  }

  function newRootScope() {
    return new Scope(null);
  }

  module.exports = {
    Scope,
    FileHash,
    newRootScope
  };
})();
