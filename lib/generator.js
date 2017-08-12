;(function() {
  "use strict";

  const _ = require("lodash");
  const uuid = require("hyperid")();
  const faker = require("faker");
  const dateformat = require("dateformat");

  function Generator(name, parent) {
    this.name = name;
    this.fields = {};
    this.fields["$id"] = new UuidField();

    if (parent) {
      this.base = parent.type();

      this.fields["$type"] = new LiteralField({value: this.type()});
      this.fields["$species"] = new LiteralField({value: this.name});
      this.fields["$extends"] = new LiteralField({value: this.base});

      _.each(Object.keys(parent.fields), (key) => {
        if (!/^\$/.test(key)) {
          this.fields[key] = new ReferenceField({key: key, generator: parent});
        }
      });
    }
  }

  Generator.prototype.type = function type() {
    if ((this.name.startsWith("$") || "" === this.name) && "" !== this.base) {
      return this.base;
    }

    return this.name;
  }

  Generator.prototype.generate = function generate(count) {
    if (count === 1) return this.one();

    var result = new Array(count);

    for (var i = 0; i < count; ++i) {
      result[i] = this.one();
    }

    return result;
  }

  Generator.prototype.one = function generateSingle() {
    return _.reduce(this.fields, (result, field, name) => {
      if (!result.hasOwnProperty(name)) {
        result[name] = field.value();
        if ("entity" === field.type) result[name]["$parent"] = result["$id"];
      }
      return result;
    }, {"$id": this.fields["$id"].value()});
  }

  Generator.prototype.withField = function resolveField(name, fieldType, options) {
    var mkField;

    switch (fieldType) {
      case "string":
        mkField = StringField;
        break;
      case "integer":
        mkField = IntegerField;
        break;
      case "decimal":
        mkField = FloatField;
        break;
      case "date":
        mkField = DateField;
        break;
      case "bool":
        mkField = BoolField;
        break;
      case "dict":
        mkField = DictField;
        break;
      case "entity":
        mkField = EntityField;
        break;
      case "literal":
        mkField = LiteralField;
        break;
      default:
        throw new Error(`Don't know how to handle ${fieldType}`);
    }

    this.fields[name] = new mkField(options);
    return this;
  };

  function Field(config) {
    this.one = function missingImpl() { throw new Error("one() must be implemented by subclasses"); }
    this.value = function generateValue() {
      if (!config || !config.countRange) {
        return this.one();
      }
      var result = [];
      for (var i = 0, count = config.countRange.count(); i < count; ++i) {
        result.push(this.one());
      }
      return result;
    }
  }

  function ReferenceField(config) {
    Field.call(this, config);
    var resolved = config.generator.fields[config.key];
    this.type = resolved.type;
    this.one = function resolveValueFromParent() { return resolved.value(); };
  }

  function EntityField(config) {
    Field.call(this, config);
    this.type = "entity";
    this.one = function nestedEntity() {
      return config.entity.one();
    };
  }

  function UuidField() {
    Field.call(this);
    this.type = "id";
    this.one = uuid;
  }

  function BoolField(config) {
    Field.call(this, config);
    this.type = "bool";
    this.one = function randBool() {
      return Math.random() > 0.49;
    };
  }

  function LiteralField(config) {
    Field.call(this, config);
    this.type = "literal";
    this.one = function constantVal() {
      return config.value;
    };
  }

  function StringField(config) {
    Field.call(this, config);
    var len = config.len, trim = 2 + (len % 14), iters = Math.floor(len / 14);

    this.type = "string";
    this.one = function randString() {
      var result = Math.random().toString(36).substring(2, trim);
      if (!iters) return result;

      for (var i = 0; i < iters; ++i) {
        result += Math.random().toString(36).substring(2, 16);
      }

      return result;
    };
  }

  function IntegerField(config) {
    Field.call(this, config);
    this.type = "integer";
    this.one = function randInt() {
      return faker.random.number(config);
    };
  }

  function FloatField(config) {
    Field.call(this, config);
    this.type = "float";
    this.one = function randFloat() {
      return parseFloat(faker.finance.amount(config.min, config.max, config.precision));
    };
  }

  function DateField(config) {
    Field.call(this, config);
    var min = dateformat(config.min, "isoUtcDateTime"), max = dateformat(config.max, "isoUtcDateTime");

    this.type = "date";
    this.one = function randDateBetween() {
      return faker.date.between(min, max);
    };
  }

  function DictField(config) {
    Field.call(this, config);

    var provider = faker, fn = config.name;

    if (fn.indexOf(".") !== -1) {
      var keys = fn.split("."), len;
      fn = keys.pop();
      len = keys.length;

      for (var i = 0, cur; (i < len) && !!(cur = keys[i]); ++i) {
        if (!provider.hasOwnProperty(cur)) unknownDictError(cur, config.name);
        provider = provider[cur];
      }
    }

    if (!provider.hasOwnProperty(fn)) unknownDictError(fn, config.name);

    this.type = "dict";
    this.one = function dictionaryChoice() {
      return provider[fn]();
    };
  }

  function unknownDictError(subName, fullName) {
    throw ((subName === fullName) ? new Error(`Cannot resolve dictionary ${subName}`) : new Error(`Cannot resolve dictionary namespace ${subName} in key ${fullName}`));
  }

  module.exports = {
    Generator
  };
})();
