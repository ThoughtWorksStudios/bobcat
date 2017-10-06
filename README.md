# Bobcat
[![wercker status](https://app.wercker.com/status/98be3d80966b1a3a006c0465c76aa8ef/s/master "wercker status")](https://app.wercker.com/project/byKey/98be3d80966b1a3a006c0465c76aa8ef)

Bobcat is a data generation tool that allows you to generate production-like data using a simple DSL. Define concepts (i.e. objects) found in your software system in our input file format, and the tool will generate JSON objects that can be inserted into a variety of datastores.

Current features include:

* Concise syntax for [modeling](#defining-entities) domain objects.
* Flexible [expression](#field-declarations) composition to generate of a variety of data.
* Over 30 built-in [dictionaries](docs/dict-basics.md) plus support for any custom dictionary to provide more realistic data values.
* [Distributions](docs/distributions.md) to determine the shape of the generated data.
* [Variable assignment](#defining-variables) for easy reference in the input file to previously generated entities.
* Ability to denote a field as the [primary key](docs/pk.md) to allow for easy insertion into a SQL database.
* File [imports](#import-statements) for better organization of input file(s).

## Table of Contents
* [Getting Started](#getting-started)
  - [User Quickstart](#user-quickstart)
  - [Developer Quickstart](#developer-quickstart)
* [Input File Format](#input-file-format)
  - [Literal Values](#literal-values)
  - [Defining Variables](#defining-variables)
  - [Defining Functions](#defining-functions)
  - [Defining Entities](#defining-entities)
  - [Generating Entities](#generating-entities)
  - [Import Statements](#import-statements)

## Getting Started

There are no prerequisites. The executable is a static binary. For more information on the usage of the executable use the flag ```--help```.

### User Quickstart

1. Download the latest [release](https://github.com/ThoughtWorksStudios/bobcat/releases)
2. Run the executable corresponding to your operating system on the sample input file:

    * Linux: `./bobcat-linux examples/example.lang`
    * macOS: `./bobcat-darwin examples/example.lang`
    * Windows: `.\bobcat-windows examples\example.lang`

3. Modify the sample file or create one from scratch to generate your own custom entities

### Developer Quickstart

1. Checkout the code:
    ```bash
    git clone https://github.com/ThoughtWorksStudios/bobcat.git
    ```
2. Set up, [build](docs/build.md), and test:
    ```bash
    make local
    ```

## Input File Format

The input file is made of three main concepts:
  * [defining entities](#defining-entities) (the objects or concepts found in your software system)
  * [fields](#field-declarations) on those entities (properties that an entity posses)
  * [generate statements](#generating-entities) to produce the desired number of entities in the resulting JSON output

The following is an example of an input file.

```example-success
# import another input file
import "examples/users.lang"

# override default $id primary key
pk("ID", $incr)

# define entity
entity Profile {
  #define fields on entity
  firstName:      $dict("first_names"),
  lastName:       $dict("last_names"),
  email:          firstName + "." + lastName + "@fastmail.com",
  addresses:      $dict("full_address")<0..3>,
  gender:         $dict("genders"),
  dob:            $date(1970-01-01, 1999-12-31, "%Y-%m-%d"),
  emailConfirmed: $bool(),
}

# declare and assign variables
let bestSelling = "Skinny"
let jeanStyles = ["Classic", "Fitted", "Relaxed", bestSelling]

entity CatalogItem {
  title: $dict("words"),
  style: $enum(jeanStyles),
  sku:   $str(10),
  price: $float(1.0, 30.00)
}

# generate statement to create corresponding JSON output
let Products = generate(10, CatalogItem)

entity CartItem {
  product:  $enum(Products),
  quantity: $int(1, 3),
}

entity Cart {
  pk("Cart_Id", $incr)
  items: CartItem<0..10>,
}

# define entity that extends an existing entity
entity Customer << User {
  last_login: $date(2010-01-01, NOW),
  profile:    Profile,
  cart:       Cart
}

# supports anonymous/inlined extensions as well
generate (10, Customer << {cart: null}) # new users don't have a cart yet
generate (90, Customer)
```

### Literal Values

| Type                           | Example                     |
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

If you need to customize the JSON representation of a literal date, you have 2 options:

1. Use `$date(min, max, format)` where `min == max` and provide a `strftime` format, e.g. `$date(2017-01-01, 2017-01-01, "%b %d, %Y")`
2. Use a literal string that looks like a date, as JSON serializes all dates as strings anyway

### Defining Variables

Declare variables with the `let` keyword followed by an [identifier](#identifiers):

```example-success
let max_value = 100
```

One does not need to initialize a declaration:

```example-success
# simply declares, but does not assign value
let foo
```

Assignment syntax should be familiar. This assigns a new value to a previous declaration:

```example-success
let max_value = 10

# assigns a new value to max_value
max_value = 1000
```

One can only assign values to variables that have been declared (i.e. implicit declarations are not supported):

```example-fail
baz = "hello" # throws error because baz was not previously declared
```

#### Identifiers

An identifier starts with a letter or underscore, followed by any number of letters, numbers, and underscores. Other symbols are not allowed. This applies to all identifiers, not just variables.

#### Predefined Variables

The following variables may be used without declaration:

| Name         | Value                                             |
|--------------|---------------------------------------------------|
| `UNIX_EPOCH` | DateTime representing `Jan 01, 1970 00:00:00 UTC` |
| `NOW`        | Current DateTime at the start of the process      |

### Defining Functions

Functions are declared using the `lambda` keyword followed by an [identifier](#identifiers), a list of input arguments, and the function body surrounded by curly braces `{}`. Note that the result of the last expression in the function body will be the return value of the function:

```example-success
# declaring perc function
lambda perc(amount, rate) {
  amount * rate
}

lambda calcTax(amount) {
  perc(amount, 0.085)
}

entity Invoice {
  price: $float(10, 30),

  # calling calcTax function on price
  tax: calcTax(price),

  total: price + tax
}
```

You can also create anonymous functions by omitting the [identifier](#identifiers) in the declaration:

```example-success
let taxRate = 0.085

entity Invoice {
  price: $float(10, 30),

  #defining anonymous function and calling it on price
  tax: (lambda (amount) { amount * taxRate })(price),

  total: price + tax
}
```

#### Native Functions

The following are functions builtin to bobcat to allow easy generation of random values. The function names are prefixed with `$` to indicate they are native, and cannot be overridden.

| Function                  | Returns                                       | Arguments                          | Defaults when omitted |
|---------------------------|-----------------------------------------------|------------------------------------|-----------------------|
| `$str(length)`            | a random string of specified length           | `length` is integer                | length=5              |
| `$float(min, max)`        | a random floating point within a given range  | `min` and `max` are numeric        | min=1.0, max=10.0     |
| `$int(min, max)`          | a random integer within a given range         | `min` and `max` are integers       | min=1, max=10         |
| `$uniqint()`          | an unsigned unique integer         | none       | none         |
| `$bool()`                 | true or false                                 | none                               | none                  |
| `$incr(offset)`           | an auto-incrementing integer from offset      | `offset` is a non-negative integer | offset=0              |
| `$uid()`                  | a 20-character unique id (MongoID compatible) | none                               | none                  |
| `$date(min, max, format)` | a random datetime within a given range        | `min` and `max` are datetimes, `format` is a `strftime` string | min=UNIX_EPOCH, max=NOW, format="%Y-%m-%dT%H:%M:%S%z" |
| `$dict(dictionary_name)`  | an entry from a specified dictionary (see [Dictionary Basics](docs/dict-basics.md) and [Custom Dictionaries](docs/dict-custom.md) for more details) | `dictionary_name` is a string | none |
| `$enum(collection)`       | [a random value from the given collection](docs/builtins.md#enumerated-field-enum) | `collection` is a collection | none |

**Note that a key difference between native functions and user-defined functions is that native functions may have optional arguments with default values.**

### Defining Entities

Entities are declared using the `entity` keyword followed by a name ([identifier](#identifiers)) and a list of [field declarations](#field-declarations) surrounded by curly braces `{}`. The entity name will be emitted as the `$type` property when the entity is serialized to JSON.

#### Field Declarations

A field declaration is simply an [identifier](#identifiers), followed by a colon `:`, an expression, and an optional [count range](docs/multi-value.md). Multiple field declarations are delimited by commas `,`. Example:

```example-success
entity User {
  # randomly selects a value from the 'email_address' dictionary
  login: $dict("email_address"),

  # creates a 16-char random-char string
  password: $str(16),

  # chooses one of the values in the collection
  status: $enum(["enabled", "disabled", "pending"])
}
```

The expressions used when defining fields can be made up of any combination of functions, literals, or references to other variables (including other fields). Right now the arithmetic operators `+ - * /` are supported.

```example-success
lambda userId(fn, ln) {
  fn + "." + ln "_" + $uniqint()
}

entity User {
  first_name: $dict("first_names"),
  last_name: $dict("last_names"),

  #compose guaranteed unique email
  email: userId(first_name, last_name) + "@" + $dict("companies") + ".com"
}
```

##### Field Distributions

To control the probability distribution of values for a specific field you can use [distributions](docs/distributions.md).

#### Anonymous Entities

Anonymous entities can be defined by omitting the identifier:

```example-success
entity {
  login: $dict("email_address"),
  status: $enum(["enabled", "disabled", "pending"])
}
```

One can also assign an anonymous entity to a variable. This allows one to reference the entity, but does not set `$type` to the variable name.

```example-success
let User = entity {
  login: $dict("email_address"),
  status: $enum(["enabled", "disabled", "pending"])
}
```

The following entity expressions are subtly different:

```example-success
# anonymous entity literal, with assignment
let Foo = entity { name: "foo" }

# formal declaration will set the entity name, as reported in the output as the `$type` property
entity Foo { name: "foo" }
```

#### Extending Entities (inheritance)

This extends the `User` entity with a `superuser` field (always set to true) into a new entity called `Admin`, whose `$type` is set to `Admin`. The original `User` entity is not modified:

```example-parse-only
entity Admin << User {
  superuser: true
}
```

As with defining other entities, extensions can be anonymous. The original User definition is not modified, and the resultant entity from the anonymous extension still reports its `$type` as `User` (i.e. the parent):

```example-parse-only
User << {
  superuser: true
}

# anonymous extension assigned to a variable
let Admin = User << {
  superuser: true
}
```

#### Entities as Fields

Field values can also be other entities:

```example-success
entity Kitten {
  says: "meow"
}

entity Person {
  name: "frank frankleton",
  pet:  Kitten
}
```

And of course any of the variations on entity expressions or declarations can be inlined here as well (see section below for more detail):

```example-success
entity Kitten {
  says: "meow"
}

entity Person {
  name:        "frank frankleton",

  # anonymous entity
  some_animal: entity { says: "oink" },

  # extended entity
  big_cat: entity Tiger << Kitten { says: "roar!" },

  # anonymous extended entity, $type is still "Kitten"
  pet:         Kitten << { says: "woof?" },
}
```

Entity fields support [multi-value](docs/multi-value.md) fields.

### Generating Entities

Ultimately one would want to generate JSON output based on the entities defined in the input file.

Generating entities is achieved with `generate(count, <entity-expression>)` statements. The default output file for the resulting JSON objects is entities.json. The entity passed in as the second argument may be defined beforehand, or inlined. `generate()` expressions return a **collection of `$id` values from each generated entity result**.

Generating 10 `User` entities:

```example-parse-only
generate(10, User) # returns a collection of the 10 `$id`s from the User entities generated
```

With anonymous entities:

```example-success
generate(10, entity {
  login: $dict("email_address"),
  password: $str(10)
})
```

Or inlined extension:

```example-parse-only
generate(10, User << {
  superuser: true
})
```

Or formally declared entities:

```example-parse-only
generate(10, entity Admin << User {
  group: "admins",
  superuser: true
})
```

### Import Statements

It's useful to organize your code into separate files for complex projects. To import other `*.lang` files, just use an import statement. Paths can be absolute, or relative to the current file:

```example-parse-only
import "path/to/file.lang"
```

