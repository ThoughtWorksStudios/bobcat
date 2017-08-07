;(function() {
  "use strict";

  const _ = require("lodash");
  const uuid = require("hyperid")();
  const faker = require("faker");
  const dateformat = require("dateformat");

  const UNIX_EPOCH = new Date(0);
  const NOW = new Date();

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
      result[name] = field.value();
      return result;
    }, {});
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

  function ReferenceField(config) {
    this.value = function resolveValueFromParent() { return config.generator.fields[config.key].value(); };
  }

  function EntityField(config) {
    this.value = function makeNested() {
      return config.entity.generate(1); // stub count to be always 1 for now
    };
  }

  function UuidField() { this.value = uuid; }

  function LiteralField(config) {
    this.value = function constantVal() {
      return config.value;
    };
  }

  function StringField(config) {
    var len = config.len
    this.value = function randString() {
      return faker.random.alphaNumeric(len);
    };
  }

  function IntegerField(config) {
    var options = _.pick(config, ["min", "max"]);

    this.value = function randInt() {
      return faker.random.number(options);
    };
  }

  function FloatField(config) {
    var options = _.pick(config, ["min", "max"]);
    options.precision = 4;

    this.value = function randFloat() {
      return faker.random.number(options);
    };
  }

  function DateField(config) {
    var min = dateformat(config.min, "isoUtcDateTime"),
      max = dateformat(config.max, "isoUtcDateTime");

    this.value = function stubString() {
      return faker.date.between(min, max);
    };
  }

  function DictField(config) {
    this.value = function stubString() {
      return `from dictionary ${config.name}`;
    };
  }

  module.exports = {
    Generator
  };
})();
