;(function() {
  "use strict";

  const   _ = require("lodash"),
         fs = require("fs"),
       Path = require("path"),
        dsl = require("./dsl"),
          s = require("./scope"),
          g = require("./generator");

  const UNIX_EPOCH = new Date(0);
  const NOW = new Date();

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

    options.countRange = parseCountRange(fieldNode.countRange);

    switch (fieldType.kind) {
      case "identifier":
        expectsArgs(0, null, "entity", fieldNode.args);
        options.entity = scope.resolve(fieldType.value);
        entity.withField(fieldNode.name, "entity", options);
        break;
      case "entity":
        expectsArgs(0, null, "entity", fieldNode.args);
        options.entity = this.entityFromNode(fieldType, scope);
        entity.withField(fieldNode.name, "entity", options);
        break;
      case "builtin":
        Object.assign(options, parseFieldArgs(fieldType.value, fieldNode.args));
        entity.withField(fieldNode.name, fieldType.value, options);
        break;
      default:
        throw new Error(`Unexpected field type ${fieldNode.kind}; field declarations must be either a built-in type or a literal value`);
    }
  };

  Interpreter.prototype.staticField = function staticField(entity, fieldNode) {
    var options = {countRange: parseCountRange(fieldNode.countRange)};
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

    var count = genNode.args[0].value, key = entity.type();

    if (this.output.hasOwnProperty(key)) {
      this.output[key] = this.output[key].concat(entity.generate(count));
    } else {
      this.output[key] = entity.generate(count);
    }
  };

  function parseCountRange(nodeSet) {
    if (!nodeSet) return null;
    if (nodeSet.length > 2) throw new Error("count range must be no more than 2 non-negative integers");

    _.each(nodeSet, (node) => {
      assertNonNegativeInt(node.value);
    });

    var min, max;

    switch (nodeSet.length) {
      case 0:
        min = max = 0;
        break;
      case 1:
        min = max = nodeSet[0].value
        break;
      default:
        min = nodeSet[0].value, max = nodeSet[1].value;
        break;
    }

    if (max < min) throw new Error("count range max must not be less than min");

    return new CountRange(min, max);
  }

  function CountRange(min, max) {
    this.count = function chooseCount() {
      if (min === max) return max;
      return _.random(min, max);
    };
  }

  function assertNonNegativeInt(number) {
    if ("number" !== typeof number || number !== Math.floor(number) || number < 0) {
      throw new Error("Expected ${JSON.stringify(number)} to be a non-negative integer");
    }
  }

  function parseFieldArgs(builtinType, nodeSet) {
    if (0 === nodeSet.length) {
      return defaultArgumentFor(builtinType)
    }

    switch (builtinType) {
      case "string":
        expectsArgs(1, "literal-int", builtinType, nodeSet);
        return {len: nodeSet[0].value};
      case "integer":
        expectsArgs(2, "literal-int", builtinType, nodeSet);
        return {min: nodeSet[0].value, max: nodeSet[1].value};
      case "decimal":
        expectsArgs(2, "literal-float", builtinType, nodeSet);
        return {min: nodeSet[0].value, max: nodeSet[1].value, precision: 4};
      case "date":
        expectsArgs(2, "literal-date", builtinType, nodeSet);
        return {min: nodeSet[0].value, max: nodeSet[1].value};
      case "bool":
        expectsArgs(0, null, builtinType, nodeSet);
        return {};
      case "dict":
        expectsArgs(1, "literal-string", builtinType, nodeSet);
        return {name: nodeSet[0].value}
      default:
        throw new Error(`Don't know how to parse args for field type ${builtinType}`);
    }
  }

  function defaultArgumentFor(builtinType) {
    switch (builtinType) {
      case "string":
        return {len: 5};
      case "integer":
        return {min: 1, max: 10};
      case "decimal":
       return {min: 1.0, max: 10.0, precision: 4};
      case "date":
       return {min: UNIX_EPOCH, max: NOW};
      case "bool":
        return {}
      default:
        throw new Error(`Field of type \`${builtinType}\` requires arguments`);
    }
  }

  function expectsArgs(num, type, fieldType, args) {
    if (!args) args = [];

    var len = args.length;

    if (num !== len) {
      throw new Error(`Field type ${fieldType} expects ${num} args, but got ${len}`);
    }

    for (var i = 0, cur, curType; i < len; ++i) {
      cur = args[i];

      if (type !== cur.kind) {
        throw new Error(`Field type ${fieldType} expects arg ${JSON.stringify(cur.value)} to be of type ${type}, but was ${cur.kind}`)
      }
    }
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
