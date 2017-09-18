# Bobcat
[![wercker status](https://app.wercker.com/status/98be3d80966b1a3a006c0465c76aa8ef/s/master "wercker status")](https://app.wercker.com/project/byKey/98be3d80966b1a3a006c0465c76aa8ef)

A data generation tool. Just define concepts in our input file format, and the tool will generate JSON objects that can be used to insert realistic-looking data into your application.

## Getting Started

### User Quickstart

1. Download the latest [release](https://github.com/ThoughtWorksStudios/bobcat/releases)
2. Run the tool over the sample input file:
    ```
    ./bobcat -o my_data.json examples/example.lang
    ```
3. Modify the sample file or create one from scratch to generate your own custom entities

### Developer Quickstart

1. Checkout the code:
    ```
    git clone https://github.com/ThoughtWorksStudios/bobcat.git
    ```
2. Set up, build, and test:
    ```
    make local
    ```

### Executable
```
Usage: bobcat [-o DESTFILE] [-d DICTPATH] [--stdout] [-cfms] [--] INPUTFILE
  bobcat -v
  bobcat -h

Arguments:
  INPUTFILE  The file describing entities and generation statements
  DESTFILE   The output file (defaults to "entities.json"); accepts "-" to output to STDOUT
  DICTPATH   The path to your user-defined dictionaries

Options:
  -h --help
  -v --version
  -c --check                           Check syntax of INPUTFILE
  -m --no-metadata                     Omit metadata in generated entities (e.g. $type, $extends, etc.)
  -o DESTFILE --output=DESTFILE        Specify output file [default: entities.json] (use "-" for DESTFILE
                                         to specify STDOUT)
  -d DICTPATH --dictionaries=DICTPATH  Specify DICTPATH
  -f --flatten                         Flattens entity hierarchies into a flat array; entities are
                                         outputted in reverse order of dependency, and linked by "$id"
  -s --split-output                    Aggregates entities by type into separate files; DESTFILE
                                         serves as the filename template, meaning each file has the
                                         entity type appended to its basename (i.e. before the ".json"
                                         extension, as in "entities-myType.json"). Implies --flatten.
  --stdout                             Alias for '-o -'; forcefully redirects output to STDOUT and
                                         supercedes setting DESTFILE elsewhere. Not compatible
                                         with --split-output.
```

### Input File Format

```
import "path/to/otherfile.lang"

entity Mammal {
  warm_blooded: true,
  says: "moo?"
}

# define entity that extends an existing entity
entity Person << Mammal {
  name:     dict("full_names"),
  roommate: Mammal { says "..." },
  pet:      entity Dog << Mammal {
    name: dict("first_names"),
    says: "oink"
  },
  login:    string(4) unique,
  dob:      date(1985-01-02, 2000-01-01),
  weight:   decimal(100.0, 250.5),
  age:      integer(21, 55),
  status:   "working",
  says:     "Greetings!"
}

generate(1, Mammal)
generate(10, Person)

# supports anonymous/inlined extensions as well
generate(5, Person << { says: "Hey you!" })
```

The input file contains definitions of entities (the objects, or concepts found in your software system), fields on those
entities (properties that an entity posses), and a 'generate' keyword to
produce the desired number of entities in the resulting JSON output. An entity has an arbitrary name,
as do fields. Entities may be nested either inline or by reference. The only other concept in this system is that of
a dictionary, which is used to provide realistic values for fields that would otherwise be difficult to generate data
for (like a person's name).

#### Import Statements

It's useful to organize your code into separate files for complex projects. To import other `*.lang` files, just use an import statement. Paths can be absolute, or relative to the current file:

```
import "path/to/file.lang"
```

#### Defining Entities

Entities are defined by curly braces that wrap a set of field definitions. For instance, this defines an anonymous entity with a login field, populated by a random email address, and a password field, populated by a 10-character random string.

##### Entity Literals

```
entity {
  login: dict("email_address"),
  password: string(10),
  status: enum(["enabled", "disabled", "pending"])
}
```

One can also simply declare a variable and assign it an anonymous entity. This allows one to reference the entity, but does not give the entity a real name as a formal entity declaration would.

```
let User = entity {
  login: dict("email_address"),
  password: string(10)
  status: enum(["enabled", "disabled", "pending"])
}
```

This also works with the entity extension syntax:

```
let Admin = User << {
  superuser: true
}
```

##### Entity Declarations
However, it's often much more useful to do an entity declaration, which sets the name of the entity; not only does this allow one to reference it later, but this **also sets the entity name** (which is reported by the `$type` property in the generated output). To formally declare an entity, provide a name (i.e. identifier) immediately after the `entity` keyword:

```
entity User {
  login: dict("email_address"),
  password: string(10)
  status: enum(["enabled", "disabled", "pending"])
}
```

The following entity expressions are subtly different:

```
# anonymous entity literal, with assignment
let Foo = entity { name: "foo" }

# formal declaration will set the entity name, as reported in the output as the `$type` property
entity Foo { name: "foo" }
```

##### Extending Entities (inheritance)

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

#### Declaring and Assigning Variables

Declare variables with the `let` keyword:

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

#### Predefined Variables

The following variables may be used without declaration:

| Name         | Value                                             |
|--------------|---------------------------------------------------|
| `UNIX_EPOCH` | DateTime representing `Jan 01, 1970 00:00:00 UTC` |
| `NOW`        | Current DateTime at the start of the process      |

#### Defining Fields

Very simply, an identifier, followed by a colon `:`, field-type, and optional arguments and count. Field declarations are delimited by commas `,`. Example:

```
entity {
  password: string(16), # creates a 16-char random-char string
  emails: dict("email_address")<1..3> # a set of 1 - 3 email addresses
}
```

##### Multi-value Field Syntax

Note that one can specify a "count range" to indicate that a field should produce an array of 0 or more values. The count range syntax is a range (lower-bound-number, followed by `..`, followed by upper-bound-number), surrounded by angled brackets (`<`, `>`).

```
# the `emails` field will yield an array of 0 - 5 email addresses.
# count ranges can be used with any field.
entity {
  emails: dict("email_address")<0..5>
}
```

Field types may be:

* One of the built-in (primitive) field types
* A literal value (for constant values)
* Another entity (inline or identifier reference)

##### Identifiers

An identifier starts with a letter or underscore, followed by any number of letters, numbers, and underscores. This applies to all identifiers, not just field names.

##### Built-in Field Types

| name            | generates                                         | arguments=(defaults)                         |
|-----------------|---------------------------------------------------|----------------------------------------------|
| string          | a string of random characters of specified length | (length=5)                                   |
| decimal         | a random floating point within a given range      | (min=1.0, max=10.0)                          |
| integer         | a random integer within a given range             | (min=1, max=10)                              |
| bool            | true or false                                     | none                                         |
| serial          | an auto-incrementing integer, starting at 1       | none                                         |
| uid             | a 12-character unique id                          | none                                         |
| date            | a date within a given range                       | (min=UNIX_EPOCH, max=NOW, optionalformat="") |
| dict            | an entry from a specified dictionary (see [Dictionary Basics](https://github.com/ThoughtWorksStudios/bobcat/wiki/Dictionary-Field-Type-Basics) and [Custom Dictionaries](https://github.com/ThoughtWorksStudios/bobcat/wiki/Creating-Custom-Dictionaries) for more details) | ("dictionary_name") -- no default |
| enum            | a random value from the given collection          | ([val1, ..., valN])                          |
| distribution    | data distribution for specified field             | (distType, fields,...)                       |

##### Literal Field Types

| type                           | example                     |
|--------------------------------|-----------------------------|
| string                         | `"hello world!"`            |
| integer                        | `1234`                      |
| decimal                        | `5.2`                       |
| bool                           | `true`                      |
| null                           | `null`                      |
| date                           | `2017-07-04`                |
| date with time                 | `2017-07-04T12:30:28`       |
| date with time (UTC)           | `2017-07-04T12:30:28Z`      |
| date with time and zone offset | `2017-07-04T12:30:28Z-0800` |
| collection (heteregenous)      | `["a", "b", "c", 1, 2, 3]`  |


##### Customizing date formats

Date fields (i.e. `date(min, max, format)`) can take an optional 3rd argument: a `strftime` format string, e.g. `"%b %d, %Y %H:%M:%S"`

If you need to customize the format of a constant date value, you have 2 options:

1. Use `date()` where min and max are the same: `date(2017-01-01, 2017-01-01, "%b %d, %Y")`
2. Use a literal string field instead, as JSON doesn't really have date types anyway (dates are always serialized to strings)

##### Entity Field Types

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

##### Enumerated Field Types (i.e. `enum`)

Enumerated values are sort of like inlined dictionaries. `enum(collection)` picks a value from the given collection:

```
# declare a collection

let statuses = ["To do", "Doing", "Done!"]

entity Work {
  status enum(statuses) # randomly picks from statuses
}
```

`generate()` statements also yield collections of `$id`s from generated entities. This can be used in conjunction with `enum` fields to relationships:

```
entity CatalogItem {
  name: string,
  sku: integer(1000, 3000)
}

# Assign the collection from generate() to a variable
let Catalog = generate(20, CatalogItem)

# each cart will have 1 - 5 CatalogItems as its contents
entity ShoppingCart {
  contents: enum(Catalog)<1..5>
}

```
##### Unique Value Flag
You can constrain the generated values for certain fields to be unique using the unique flag. The following is an example using the unique flag.

```
entity CatelogItem {
  name: string(10) unique,
  sku:  integer(1000, 3000)
}
```

It's important to note that boolean, static, and entity field types don't support the unique flag, and that it may not be possible to provide unique values under certain conditions. The following example is a case where there don't exist enough unique possible values which will cause an error to be returned.

```
entity Human {
  name: dict("full_names"),
  age:  integer(1, 10)
}

generate(50, Human)
```

Since there are only 10 possible values for the age field, it's not possible to generate 50 Humans with each age value being unique.

#### Generating Entities (i.e. Generate Expressions)

Generating entities is achieved with `generate(count, <entity-expression>)` statements. The entity passed in as the second argument may be defined beforehand, or inlined. `generate()` expressions return a **collection of `$id` values from each generated entity result**.

Generating 10 `User` entities:

```
generate(10, User) # returns a collection of the 10 `$id`s from the User entities generated
```

With anonymous entities:

```
generate(10, entity {
  login: dict("email_address"),
  password: string(10)
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

#### Distributions

Distribution fields allow you the specify the shape that the generated data should take.

The following are currently supported distributions

| Name         | Value                                             | Allowed Fields   |
|--------------|---------------------------------------------------|------------------|
| `normal`     | The normal gaussian distribution                  | Decimal          |
| `uniform`    | A uniform distribution                            | integer, decimal |
| `percent`    | specify the % something should occur              | all              |
| `weighted`   | probability weights                               | all              |


example:
```
entity Human {
  age: distribution(percent,
    25% => decimal(1.0, 15.0),
    50% => decimal(15.0, 30.0),
    25% => decimal(30.0, 80.0)
  ),
  height: distribution(weighted,
    55 => decimal(1.0, 15.0),
    500 => decimal(15.0, 30.0),
    2 => decimal(30.0, 80.0)
  ),
  weight: distribution(normal, decimal(1.0, 400.0))
}
```

### Prerequisites

None. The executable is a static binary.

### Building from Source

The included Makefile has targets to get you started:

```
$ make list

Make targets:
   build clean compile depend local performance prepare release run smoke test wercker
```

Set up your dev workspace. This will install golang from homebrew, configure the current directory for development, install dependencies, then finally build and run tests:

```
make local
```

Build and run tests:
```
make
```

Just build the binary:
```
make build
```

Just run tests:
```
make test
```
