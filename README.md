# Bobcat
[![wercker status](https://app.wercker.com/status/98be3d80966b1a3a006c0465c76aa8ef/s/master "wercker status")](https://app.wercker.com/project/byKey/98be3d80966b1a3a006c0465c76aa8ef)

A data generation tool. Just define concepts in our input file format, and the tool will generate JSON objects that can be used to insert realistic-looking data into your application.

## Getting Started

### User Quickstart

1. Download the latest [release](https://github.com/ThoughtWorksStudios/bobcat/releases)
2. Run the tool over the sample input file:
    ```
    ./bobcat -dest=my_data.json examples/example.lang
    ```
3. Modify the sample file or create one from scratch to generate your own custom entities

### Developer Quickstart

1. [Install Docker for Mac](https://download.docker.com/mac/stable/Docker.dmg)
2. Checkout the code:
    ```
    git clone https://github.com/ThoughtWorksStudios/bobcat.git
    ```
3. Set up, build, and test:
    ```
    make local
    ```
4. Alternatively, start the Docker container to do the same as the previous step, but in a container:
    ```
    make docker
    ```

### Executable
```
Usage: ./bobcat [ options ] spec_file.lang

Options:
  -c
      Checks the syntax of the provided spec
  -d string
      location of custom dictionary files ( e.g. ./bobcat -d=~/data/ examples/example.lang )
  -dest string
      Destination file for generated content (NOTE that -dest and -split-output are mutually exclusize; the -dest flag will be ignored) (default "entities.json")
  -split-output
      Create a seperate output file per definition with the filename being the definition's name. (NOTE that -split-output and -dest are mutually exclusize; the -dest flag will be ignored)
```
### Input file format

```
import "path/to/otherfile.lang"

Mammal: {
  warm_blooded true,
  says "moo?"
}

Person: Mammal {
  name     dict("full_names"),
  roommate Mammal { says "..." },
  pet      Dog:Mammal {
    name dict("first_names"),
    says "oink"
  },
  login    string(4),
  dob      date(1985-01-02, 2000-01-01),
  weight   decimal(100.0, 250.5),
  age      integer(21, 55),
  status   "working",
  says     "Greetings!"
}

generate (1, Mammal)
generate (10, Person)
generate (5, Person { says "Hey you!" })
```

The input file contains definitions of entities (the objects, or concepts found in your software system), fields on those
entities (properties that an entity posses), and a 'generate' keyword to
produce the desired number of entities in the resulting JSON output. An entity has an arbitrary name,
as do fields. Entities may be nested either inline or by reference. The only other concept in this system is that of
a dictionary, which is used to provide realistic values for fields that would otherwise be difficult to generate data
for (like a person's name).

#### Defining entities

Entities are defined by curly braces that wrap a set of field definitions. For instance, this defines an anonymous entity with a login field, populated by a random email address, and a password field, populated by a 10-character random string:

```
{
  login dict("email_address"),
  password string(10),
  status enum("enabled", "disabled", "pending")
}
```

It's much more useful to name this entity so that one can reference it later. To do this, simply assign the entity definition to an identifier followed by a colon `:`:

```
User: {
  login dict("email_address"),
  password string(10)
  status enum("enabled", "disabled", "pending")
}
```

#### Defining fields

Very simply, an identifier, followed by a field-type and optional arguments. Example:

```
password string(16)
```

Field types may be:

* One of the built-in (primitive) field types
* A literal value (for constant values)
* Another entity (inline or identifier reference)

##### Identifiers

An identifier starts with a letter or underscore, followed by any number of letters, numbers, and underscores. This applies to all identifiers, not just field names.

##### Built-in field types

| name    | generates                                         | arguments=(defaults)      |
|---------|---------------------------------------------------|---------------------------|
| string  | a string of random characters of specified length | (length=5)                |
| decimal | a random floating point within a given range      | (min=1.0, max=10.0)       |
| integer | a random integer within a given range             | (min=1, max=10)           |
| bool    | true or false                                     | none                      |
| date    | a date within a given range                       | (min=UNIX_EPOCH, max=NOW) |
| dict    | an entry from a specified dictionary (see [Dictionary Basics](https://github.com/ThoughtWorksStudios/bobcat/wiki/Dictionary-Field-Type-Basics) and [Custom Dictionaries](https://github.com/ThoughtWorksStudios/bobcat/wiki/Creating-Custom-Dictionaries) for more details) | ("dictionary_name") -- no default |
| enum    | One of the provided values                       | (string, string, ...) |

##### Literal types

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

##### Entity types

Entity types can be declared by just referencing the identifier:

```
Kitteh: {
  says "meh"
}

Person: {
  name "frank frankleton",
  pet  Kitteh
}
```

And of course any of the variations on entity declarations can be inlined here as well (see section below for more detail):

```
Kitteh: {
  says "meh"
}

Person: {
  name        "frank frankleton",
  pet         Kitteh { says "meow?" },
  some_animal { says "oink" }
}
```

#### Extending entities (inheritance)

This extends the `User` entity with a `superuser` field (always set to true) into a new entity called `Admin`. The original `User` entity is not modified:

```
Admin: User {
  superuser true
}
```

As with defining other entities, one does not have to assign an identifier and can instead do this as a one-off "anonymous extension." The original User definition is not modified:

```
User {
  superuser true
}
```

#### Import statements

It's useful to organize your code into separate files for complex projects. To import other `*.lang` files, just use an import statement. Paths can be absolute, or relative to the current file:

```
import "path/to/file.lang"
```

#### Generating entities

Generating entities is achieved with `generate(count, entity)` statements. The entity passed in as the second argument may be defined beforehand, or inlined.

Generating 10 `User` entities:

```
generate(10, User)
```

With anonymous entities:

```
generate(10, {
  login dict("email_address"),
  password string(10)
})
```

Or inlined extension:

```
generate(10, User {
  superuser true
})
```

Or inlined, named extension:

```
generate(10, Admin:User {
  group: "admins",
  superuser true
})
```

### Prerequisites

None. The executable is a static binary.

### Building from source

The included Makefile has targets to get you started:

```
make list
  build clean depend docker local release run test wercker
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

Alternatively, one can use Docker to develop, build, and run. The simplest way to do this is to install Docker for Mac and then run:

```
make docker
```

This will create a docker container, build the software, and run the example file.
