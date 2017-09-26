# Bobcat
[![wercker status](https://app.wercker.com/status/98be3d80966b1a3a006c0465c76aa8ef/s/master "wercker status")](https://app.wercker.com/project/byKey/98be3d80966b1a3a006c0465c76aa8ef)

Bobcat is a data generation tool that allows you to generate production-like data using a simple DSL. Define concepts (i.e. objects) found in your software system in our input file format, and the tool will generate JSON objects that can be inserted into a variety of datastores.

Current features include:

* Concise syntax for [modeling](#defining-entities) domain objects.
* Flexible [field types](#defining-fields) for generation of a variety of data.
* Over 30 built-in [dictionaries](docs/dict-basic.md) plus support for any custom dictionary to provide more realistic data values.
* [Distributions](docs/builtins.md#distribution-field) to determine the shape of the generated data.
* [Variable assignment](#declaring-and-assigning-variables) for easy reference in the input file to previously generated entities.
* The ability to denote a field as the primary key to allow for easy insertion into a SQL database.
* [Unique](#docs/builtins.md#unique-value-flag) field flag so the values generated for that field will be unique over the collection of JSON objects.
* File [imports](#import-statements) for better organization of input file(s).

## Table of Contents
* [Getting Started](#getting-started)
  - [User Quickstart](#user-quickstart)
  - [Developer Quickstart](#developer-quickstart)
* [Input File Format](#input-file-format)
  - [Defining Entities](#defining-entities)
  - [Defining Fields](#defining-fields)
  - [Generating Entities](#generating-entities-generate-expressions)
  - [Variable Assignment](#declaring-and-assigning-variables)
  - [Import Statements](#import-statements)

## Getting Started

There are no prerequisites. The executable is a static binary. For more information on the usage of the executable use the flag ```--help```.

### User Quickstart

1. Download the latest [release](https://github.com/ThoughtWorksStudios/bobcat/releases)
2. Run the executable corresponding to your operating system on the sample input file:

    Linux:
    ```
    ./bobcat-linux examples/example.lang
    ```
    macOS:
    ```
    ./bobcat-darwin examples/example.lang
    ```
    Windows:
    ```
    .\bobcat-windows examples\example.lang
    ```
3. Modify the sample file or create one from scratch to generate your own custom entities

### Developer Quickstart

1. Checkout the code:
    ```
    git clone https://github.com/ThoughtWorksStudios/bobcat.git
    ```
2. Set up, [build](docs/build.md), and test:
    ```
    make local
    ```

## Input File Format

The input file is made of three main concepts:
  * definitions of entities (the objects or concepts found in your software system)
  * fields on those entities (properties that an entity posses)
  * generate statements to produce the desired number of entities in the resulting JSON output

The input file also supports variable assignment and import statements.

The following is an example of an input file.

```
import "users.lang"

let SHORT_DATE_FORMAT = "%Y-%m-%d"

# define entity
entity Profile {
  #define fields on entity
  firstName:      $dict("first_names"),
  lastName:       $dict("last_names"),
  email:          firstName + "." + lastName + "@fastmail.com",
  addresses:      $dict("full_address")<0..3>,
  gender:         $dict("genders"),
  dob:            $date(1970-01-01, 1999-12-31, SHORT_DATE_FORMAT),
  emailConfirmed: $bool(),
}

#declare and assign variables
let bestSelling = "Skinny"
let jeanStyles = ["Classic", "Fitted", "Relaxed", bestSelling]

entity CatalogItem {
  title: $dict("words"),
  style: $enum(jeanStyles),
  sku:   $str(10),
  price: $float(1.0, 30.00)
}

#generate statement to create corresponding JSON output
let Products = generate(10, CatalogItem)

entity CartItem {
  product:  $enum(Products),
  quantity: $int(1, 3),
}

entity Cart {
  items: CartItem<0..10>,
  total: $float() # TODO: this should be a calculated item based on CartItems price x quantity + tax
}

# define entity that extends an existing entity
entity Customer << User {
  last_login: $date(2010-01-01, NOW), # UNIX_EPOCH and NOW are predefined variables
  profile:    Profile,
  cart:       Cart
}

# supports anonymous/inlined extensions as well
generate (10, Customer << {cart: null}) # new users don't have a cart yet
generate (90, Customer)
```

### Defining Entities

Entities are defined by curly braces that wrap a set of field definitions. For instance, this defines an anonymous entity with a login field (populated by a random email address), a password field (populated by a 10-character random string), and a status field (populated from the options "enabled", "disabled", "pending").

#### Entity Literals

```
entity {
  login: $dict("email_address"),
  password: $str(10),
  status: $enum(["enabled", "disabled", "pending"])
}
```

One can also simply declare a variable and assign it an anonymous entity. This allows one to reference the entity, but does not give the entity a real name as a formal entity declaration would.

```
let User = entity {
  login: $dict("email_address"),
  password: $str(10),
  status: $enum(["enabled", "disabled", "pending"])
}
```

This also works with the entity extension syntax:

```
let Admin = User << {
  superuser: true
}
```

#### Entity Declarations
However, it's often much more useful to do an entity declaration, which sets the name of the entity; not only does this allow one to reference it later, but this **also sets the entity name** (which is reported by the `$type` property in the generated output). To formally declare an entity, provide a name ([identifier](#identifiers)) immediately after the `entity` keyword:

```
entity User {
  login: $dict("email_address"),
  password: $str(10),
  status: $enum(["enabled", "disabled", "pending"])
}
```

The following entity expressions are subtly different:

```
# anonymous entity literal, with assignment
let Foo = entity { name: "foo" }

# formal declaration will set the entity name, as reported in the output as the `$type` property
entity Foo { name: "foo" }
```

#### Extending Entities (inheritance)

This extends the `User` entity with a `superuser` field (always set to true) into a new entity called `Admin`, whose `$type` is set to `Admin`. The original `User` entity is not modified:

```
entity Admin << User {
  superuser: true
}
```

As with defining other entities, one does not have to assign an identifier / formally declare a descendant entity; extension expressions can be anonymous. The original User definition is not modified, and the resultant entity from the anonymous extension still reports its `$type` as `User` (i.e. the parent):

```
User << {
  superuser: true
}
```

### Defining Fields

A Field declaration is simply an [identifier](#identifiers), followed by a colon `:`, field-type, and optional arguments and [count](docs/multi-value.md). Multiple field declarations are delimited by commas `,`. Example:

```
entity {
  # creates a 16-char random-char string
  password: $str(16),

  # the last field declaration may have a trailing comma
  emails: $dict("email_address"),
}
```

Field types may be:

* One of the built-in field types
* A literal value (for constant values)
* Another entity (inline or identifier reference)
* A calculated value

#### Built-in Field Types

| name            | generates                                         | arguments=(defaults)                         | supports [unique](docs/builtins.md#unique-value-flag) |
|-----------------|---------------------------------------------------|----------------------------------------------|----------------------|
| `$str()`        | a string of random characters of specified length | (length=5)                                   | yes                  |
| `$float()`      | a random floating point within a given range      | (min=1.0, max=10.0)                          | yes                  |
| `$int()`        | a random integer within a given range             | (min=1, max=10)                              | yes                  |
| `$bool()`       | true or false                                     | none                                         | no                   |
| `$incr()`       | an auto-incrementing integer, starting at 1       | none                                         | yes                  |
| `$uid()`        | a 20-character unique id (MongoID compatible)     | none                                         | yes                  |
| [`$date()`](docs/builtins.md#customizing-date-formats)            | a date within a given range                    | (min=UNIX_EPOCH, max=NOW, optionalformat="") | yes                  |
| `$dict()`       | an entry from a specified dictionary (see [Dictionary Basics](docs/dict-basics.md) and [Custom Dictionaries](docs/dict-custom.md) for more details) | ("dictionary_name") -- no default | yes                   |
| [`$enum()`](docs/builtins.md#enumerated-field-enum )         | a random value from the given collection          | ([val1, ..., valN])                          | yes                   |
| [`$distribution()`](docs/builtins.md#distribution-field)    | data distribution for specified field             | none                                         | no                   |

More information about built-in fields can be found [here](docs/builtins.md).

#### Literal Field Types

| type                           | example                     |
|--------------------------------|-----------------------------|
| string                         | `"hello world!"`            |
| integer                        | `1234`                      |
| float                          | `5.2`                       |
| bool                           | `true`                      |
| null                           | `null`                      |
| date                           | `2017-07-04`                |
| date with time                 | `2017-07-04T12:30:28`       |
| date with time (UTC)           | `2017-07-04T12:30:28Z`      |
| date with time and zone offset | `2017-07-04T12:30:28Z-0800` |
| collection (heteregenous)      | `["a", "b", "c", 1, 2, 3]`  |

If you need to customize the format of a literal date field (constant date value), you have 2 options:

1. Use `date()` where min and max are the same: `date(2017-01-01, 2017-01-01, "%b %d, %Y")`
2. Use a literal string field instead, as JSON doesn't really have date types anyway (dates are always serialized to strings)

#### Entity Field Types

Entity fields can be declared by just referencing an entity by identifier:

```
entity Kitteh {
  says: "meh"
}

entity Person {
  name: "frank frankleton",
  pet:  Kitteh
}
```

And of course any of the variations on entity expressions or declarations can be inlined here as well (see section below for more detail):

```
entity Kitteh {
  says: "meh"
}

entity Person {
  name:        "frank frankleton",

  # anonymous extension, $type is still "Kitteh"
  pet:         Kitteh << { says: "meow?" },

  some_animal: { says: "oink" }, # anonymous entity

  # formal declarations are expressions too
  big_cat: entity Tiger << Kitteh { says: "roar!" }
}
```

Entity fields support [multi-value](docs/multi-value.md) fields.

#### Calculated Field Types

Calculated field types can include literal values or references to other fields or variables (identifers). Right now the arithmetic operators ```+ - * / ``` are supported.

```
let tax_rate = 0.0987

entity Product {

  price: $float(1.00, 300.00),
  quantity: $int(1, 5),
  sub_total: price * quantity,
  tax: sub_total * tax_rate,
  total: sub_total + tax
}

```

### Generating Entities (Generate Expressions)

Generating entities is achieved with `generate(count, <entity-expression>)` statements. The entity passed in as the second argument may be defined beforehand, or inlined. `generate()` expressions return a **collection of `$id` values from each generated entity result**.

Generating 10 `User` entities:

```
generate(10, User) # returns a collection of the 10 `$id`s from the User entities generated
```

With anonymous entities:

```
generate(10, entity {
  login: $dict("email_address"),
  password: $str(10)
})
```

Or inlined extension:

```
generate(10, User << {
  superuser: true
})
```

Or formally declared entities:

```
generate(10, entity Admin << User {
  group: "admins",
  superuser: true
})
```

### Import Statements

It's useful to organize your code into separate files for complex projects. To import other `*.lang` files, just use an import statement. Paths can be absolute, or relative to the current file:

```
import "path/to/file.lang"
```

### Declaring and Assigning Variables

Declare variables with the `let` keyword followed by an identifier:

```
let max_value = 100
```

One does not need to initialize a declaration:

```
# simply declares, but does not assign value
let foo
```

Assignment syntax should be familiar:

```
let max_value = 10

# assigns a new value to max_value
max_value = 1000
```

One can only assign values to variables that have been declared (i.e. implicit declarations are not supported):

```
baz = "hello" # throws error as baz was not previously declared
```

#### Identifiers

An identifier starts with a letter or underscore, followed by any number of letters, numbers, and underscores. This applies to all identifiers, not just variables.

#### Predefined Variables

The following variables may be used without declaration:

| Name         | Value                                             |
|--------------|---------------------------------------------------|
| `UNIX_EPOCH` | DateTime representing `Jan 01, 1970 00:00:00 UTC` |
| `NOW`        | Current DateTime at the start of the process      |
