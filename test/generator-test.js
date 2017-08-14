const test = require("tape"),
         _ = require("lodash"),
       gen = require("../lib/generator");

const Generator = gen.Generator,
 ReferenceField = gen.ReferenceField,
    EntityField = gen.EntityField,
      UuidField = gen.UuidField,
      BoolField = gen.BoolField,
   LiteralField = gen.LiteralField,
    StringField = gen.StringField,
   IntegerField = gen.IntegerField,
     FloatField = gen.FloatField,
      DateField = gen.DateField,
      DictField = gen.DictField;

test("generator extension/inheritance", (t) => {
  let parent = new Generator("base").
    withField("name", "literal", {value: "big kyle"}).
    withField("species", "literal", {value: "h00man"}).
    withField("age", "integer", {min: 20, max: 30});

  let child = new Generator("child", parent).
    withField("name", "literal", {value: "lil kyle"}).
    withField("age", "integer", {min: 5, max: 10});

  let bigKyle = parent.one();
  let lilKyle = child.one();

  t.equal(bigKyle.name, "big kyle");
  t.equal(bigKyle.species, "h00man");
  t.ok(bigKyle.age >= 20 && bigKyle.age <= 30);

  t.equal(lilKyle.name, "lil kyle", "child overrides static field");
  t.equal(lilKyle.species, "h00man", "child should inherit parent's field");
  t.ok(lilKyle.age >= 5 && lilKyle.age <= 10, "child overrides dynamic field");

  t.end();
});

test("entity A nests entity B as an entity field ", (t) => {
  let inner = new Generator("innie").
    withField("name", "literal", {value: "tapeworm"});

  let outer = new Generator("outtie").
    withField("name", "literal", {value: "kyle"}).
    withField("parasite", "entity", {entity: inner});

  let kyle = outer.one();

  t.ok(kyle.parasite, "entity A contains entity B");
  t.equal(typeof kyle.$id, "string", "entity A should have an $id");
  t.equal(typeof kyle.parasite.$id, "string", "nested entity B also has an $id");
  t.notEqual(kyle.parasite.$id, kyle.$id, "entities A and B have different $id values");
  t.test("$parent field in nested instances", (t) => {
    t.equal(kyle.parasite.$parent, kyle.$id, "should be set to $id of outer entity");
    t.notOk(inner.one().hasOwnProperty("$parent"), "should be absent in unnested instances generated afterward; $parent should not be added to entity B's definition");
    t.end();
  });

  t.end();
});

test("withField() creates the correct field types", (t) => {
  let g = new Generator("foo");
  let specs = [
    {type: "string", expected: StringField},
    {type: "integer", expected: IntegerField},
    {type: "decimal", expected: FloatField},
    {type: "date", expected: DateField},
    {type: "bool", expected: BoolField},
    {type: "dict", expected: DictField},
    {type: "entity", expected: EntityField},
    {type: "literal", expected: LiteralField}
  ];

  _.each(specs, (sp) => { g.withField(sp.type, sp.type, {name: "name.firstName" /* satisfies dictionary */})});
  _.each(specs, (sp) => {
    t.equal(g.fields[sp.type].constructor, sp.expected, `${sp.type} should yield a ${sp.expected.name}`);
  });

  t.throws(() => {
    g.withField("foo", "foo", {});
  }, /Don't know how to handle foo/, "should throw error on unknown field type");

  t.end();
});

test("generate() produces specified output", (t) => {
  let g = new Generator("foo");
  let e = new Generator("other");
  let start = new Date("2017-01-01"), end = new Date("2017-02-01");

  let specs = [
    {type: "string", test: stringOfLen, options: {len: 10}},
    {type: "integer", test: intWithin, options: {min: 20, max: 23}},
    {type: "decimal", test: floatWithin, options: {min: 4.5, max: 6.5, precision: 4}},
    {type: "date", test: dateWithin, options: {min: start, max: end}},
    {type: "bool", test: isBool, options: {}},
    {type: "entity", test: isEntity, options: {entity: e}},
    {type: "dict", test: isString, options: {name: "name.firstName"}},
    {type: "literal", options: {value: "hi"}}
  ];

  _.each(specs, (sp) => { g.withField(sp.type, sp.type, sp.options); });

  let foo = g.one();

  t.ok(foo.hasOwnProperty("$id"), "should generate $id");
  t.equal("string", typeof foo.$id, "$id is a string");

  _.each(specs, (sp) => {
    let actual = foo[sp.type];
    if (sp.type === "literal") {
      t.equal(actual, sp.options.value, `field ${sp.type} should return value`);
    } else {
      t.ok(sp.test(actual, sp.options), `field ${sp.type} should match ${sp.test.name}(${JSON.stringify(sp.options)})`);
    }
  });
  t.end();
});

test("withField() honors count when provided", (t) => {
  let g = new Generator("foo");
  let e = new Generator("other");
  let start = new Date("2017-01-01"), end = new Date("2017-02-01");

  let specs = [
    {type: "string", options: {len: 10}, expected: 3},
    {type: "integer", options: {min: 20, max: 23}, expected: 1},
    {type: "decimal", options: {min: 4.5, max: 6.5, precision: 4}, expected: 2},
    {type: "date", options: {min: start, max: end}, expected: 5},
    {type: "bool", options: {}, expected: 9},
    {type: "entity", options: {entity: e}, expected: 2},
    {type: "dict", options: {name: "name.firstName"}, expected: 5},
    {type: "literal", options: {value: "hi"}, expected: 6}
  ];

  _.each(specs, (sp) => { g.withField(sp.type, sp.type, Object.assign({countRange: count(sp.expected)}, sp.options)); });

  let foo = g.one();

  t.equal("string", typeof foo.$id, "$id is always a single string");

  _.each(specs, (sp) => {
    t.ok(foo[sp.type] instanceof Array, `field ${sp.type} should be an array`);
    t.equal(foo[sp.type].length, sp.expected, `field ${sp.type} should honor the count param`);
  });
  t.end();
});

test("generate() honors count when provided", (t) => {
  let g = new Generator("foo");

  let result = g.generate(3);
  let unique = _.uniq(_.compact(_.map(result, (r) => {return r.$id})));
  t.equal(result.length, 3, "should generate the specified count");
  t.equal(unique.length, 3, "each entity should be unique");
  t.end();
});

function count(num) { // mock count range object
  return {count: () => {return num;}};
}

function stringOfLen(actual, expected) {
  return isString(actual) && actual.length === expected.len;
}

function intWithin(actual, expected) {
  return "number" === typeof actual && precision(actual, 0) && numBtw(actual, expected.min, expected.max);
}

function floatWithin(actual, expected) {
  return "number" === typeof actual && precision(actual, expected.precision) && numBtw(actual, expected.min, expected.max);
}

function dateWithin(actual, expected) {
  return (actual instanceof Date) && numBtw(actual.getTime(), expected.min.getTime(), expected.max.getTime());
}

function isBool(actual) {
  return "boolean" === typeof actual;
}

function isEntity(actual) {
  return "string" === typeof actual.$id;
}

function isString(actual) {
  return "string" === typeof actual;
}

function numBtw(val, min, max) {
  return val >= min && val <= max;
}

function precision(val, pre) {
  let frac = (val + "").split(".")[1];
  if (pre === 0) return !frac;
  return frac.length === pre;
}
