;(function() {
  "use strict";

  const   _ = require("lodash"),
         fs = require("fs"),
       Path = require("path"),
        dsl = require("./dsl"),
          s = require("./scope"),
          g = require("./generator");

  function Interpreter() {
    this.anonNSCounter = new NamespacedCounter();
    this.baseDir = ".";
    this.output = {};
  }

  Interpreter.prototype.loadFile = function loadFile(fspath, scope) {
    var original = this.baseDir,
        realpath = resolve(fspath, original),
        tree;

    if (scope.imports.hasSeen(realpath)) {
      console.log(`already loaded ${realpath}`);
      return;
    }

    this.baseDir = calculateNewBasedir(fspath, original);

    tree = dsl.parse(fs.readFileSync(realpath, "utf-8"));

    scope.imports.markSeen(realpath);

    this.visit(tree, scope);

    this.baseDir = original;
  };

  Interpreter.prototype.visit = function visit(node, scope) {
    switch (node.kind) {
      case "root":
        _.each(node.children, (child) => {
          this.visit(child, scope);
        });
        break;
      case "entity":
        this.entityFromNode(node, scope);
        break;
      case "generation":
        this.generateFromNode(node, scope);
        break;
      case "import":
        this.loadFile(node.value, scope);
        break;
      default:
        throw new Error(`Unexpected token type ${node.kind}`);
    }
  };

  Interpreter.prototype.entityFromNode = function entityFromNode(entNode, scope) {
    var entity, formalName, parentScope = scope;
    scope = new s.Scope(parentScope); // new child scope

    if (entNode.related) {
      let symbol = entNode.related.value;
      let parentEntity = parentScope.resolve(symbol);

      formalName = !entNode.name ? ["$" + this.anonNSCounter.next(symbol), symbol].join("::") : entNode.name;

      entity = new g.Generator(formalName, parentEntity);
    } else {
      formalName = !entNode.name ? "$" + this.anonNSCounter.next("$") : entNode.name;
      entity = new g.Generator(formalName);
    }

    // Add entity to symbol table before iterating through field defs so fields can reference
    // the current entity.
    parentScope.set(formalName, entity);

    _.each(entNode.children, (fieldNode) => {
      var fieldType = fieldNode.value;
      if (/^literal-/.test(fieldType.kind)) {
        this.staticField(entity, fieldNode);
      } else {
        this.dynamicField(entity, fieldNode, scope);
      }
    });

    return entity;
  };

  Interpreter.prototype.dynamicField = function dynamicField(entity, fieldNode, scope) {
    var fieldType = fieldNode.value;
    var sourceEntity, options = {};

    switch (fieldType.kind) {
      case "identifier":
        sourceEntity = scope.resolve(fieldType.value);
        options = { entity: sourceEntity };
        entity.withField(fieldNode.name, "entity", options);
        break;
      case "entity":
        sourceEntity = this.entityFromNode(fieldType, scope);
        options = { entity: sourceEntity };
        entity.withField(fieldNode.name, "entity", options);
        break;
      case "builtin":
        entity.withField(fieldNode.name, fieldType.value, fieldType.args || {});
        break;
      default:
        throw new Error(`Unexpected field type ${fieldNode.kind}; field declarations must be either a built-in type or a literal value`);
    }
  };

  Interpreter.prototype.staticField = function staticField(entity, fieldNode) {
    var options = detectMultiple(fieldNode.bound)
    options.value = fieldNode.value.value;
    entity.withField(fieldNode.name, "literal", options);
  };

  Interpreter.prototype.generateFromNode = function generateFromNode(genNode, scope) {
    var type = genNode.value.kind;
    var entity;

    switch (type) {
      case "identifier":
        entity = scope.resolve(genNode.value.value);
        break;
      case "entity":
        entity = this.entityFromNode(genNode.value, scope);
        break;
      default:
        throw new Error(`Don't know how to generate a(n) ${type}`);
    }

    var count = genNode.args[0].value;
    console.log(JSON.stringify(entity.generate(count), null, 2));
  };

  function detectMultiple(count) {
    if (!count) return {};
    return { multiple: true, lower: bound.min, upper: bound.max };
  }

  function resolve(fspath, base) {
    fspath = Path.normalize(fspath);

    if (Path.isAbsolute(fspath) || base === "") {
      return fs.realpathSync(fspath);
    }

    return fs.realpathSync(Path.join(base, fspath));
  }

  function calculateNewBasedir(fspath, base) {
    fspath = Path.normalize(fspath);

    if (Path.isAbsolute(fspath) || base === "") {
      return Path.dirname(fspath);
    }

    return Path.relative(".", Path.join(base, Path.dirname(fspath)));
  }

  function NamespacedCounter() {
    var names = {};
    this.next = function nextInt(namespace) {
      if (!namespace) throw new Error("Cannot increment counter without namespace!");
      if ("undefined" === typeof names[namespace]) {
        names[namespace] = 1;
      } else {
        names[namespace] += 1;
      }

      return names[namespace];
    }
  }

  module.exports = Interpreter;
})();
