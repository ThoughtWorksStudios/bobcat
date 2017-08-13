const test = require("tape");

const s = require("../lib/scope");

test("scope symbol resolution", function scopeSymbolResolution(t) {
  let scope = s.newRootScope();
  scope.set("foo", 1024);

  t.doesNotThrow(() => {
    t.equal(scope.resolve("foo"), 1024, "should resolve to the same value");
  }, /Failed to resolve symbol/, "should resolve symbols declared in current scope");

  t.test("with child scope", function resolvesFromChildScope(t) {
    let child = new s.Scope(scope);
    child.set("bar", "baz");
    t.doesNotThrow(() => {
      t.equal(child.resolve("bar"), "baz", "should resolve to same value");

      t.equal(child.resolve("foo"), 1024, "should resolve symbols declared in parent scope");
      child.set("foo", 3);
      t.equal(child.resolve("foo"), 3, "can override symbols declared in parent");
      t.equal(scope.resolve("foo"), 1024, "overriding symbols in child does not affect parent");
    }, /Failed to resolve symbol/, "should resolve symbols in child scope");

    t.throws(() => {
      scope.resolve("bar");
    }, /Failed to resolve symbol/, "parent scope should not resolve symbols declared in child scope")

    t.end();
  });
  t.end();
});
