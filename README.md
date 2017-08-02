# DataGen
[![wercker status](https://app.wercker.com/status/98be3d80966b1a3a006c0465c76aa8ef/s/master "wercker status")](https://app.wercker.com/project/byKey/98be3d80966b1a3a006c0465c76aa8ef)

A data generation tool. Just define concepts in our input file format, and the tool will generate JSON objects that can be used to insert realistic-looking data into your application.

## Getting Started

### User Quickstart

1. Download the latest [release](https://github.com/ThoughtWorksStudios/datagen/releases)
2. Run the tool over the sample input file:

        ./datagen -dest=my_data.json example.lang

3. Modify the sample file or create one from scratch to generate your own custom entities

### Developer Quickstart

1. [Install Docker for Mac](https://download.docker.com/mac/stable/Docker.dmg)
2. Checkout the code:

        git clone https://github.com/ThoughtWorksStudios/datagen.git

3. Start the Docker container:

        make docker

Note: if you prefer local development over a docker container, try 'make local'. And see 'make list' for all commands.

### Executable
```
Usage: ./datagen [ options ] spec_file.lang

Options:
  -c
      Checks the syntax of the provided spec
  -d string
      location of custom dictionary files ( e.g. ./datagen -d=~/data/ example.lang )
  -dest string
      Destination file for generated content (NOTE that -dest and -split-output are mutually exclusize; the -dest flag will be ignored) (default "entities.json")
  -split-output
      Create a seperate output file per definition with the filename being the definition's name. (NOTE that -split-output and -dest are mutually exclusize; the -dest flag will be ignored)
```
### Input file format

```
thing: {
  exists "false"
}

Person: thing {
  full_name dict("full_names"),
  login string(4),
  dob date(1985-01-02, 2000-01-01),
  age  decimal(4.2, 42.7),
  status "working",
  exists "true",
}

generate (1, thing)
generate (5, Person { status "hmmm" })
```

The input file contains definitions of entities (the objects, or concepts found in your software system), fields on those
entities (properties that an entity posses), and a 'generate' keyword to
produce the desired number of entities in the resulting JSON output. An entity has an arbitrary name,
as do fields. The only other concept in this system is that of a dictionary, which is used to provide
realistic values for fields that would otherwise be difficult to generate data for (like a person's name).

#### Defining entities

```
thing: { /* field definitions go here*/ }
```

#### Defining fields

```
login string(4),
age  decimal(4.2, 42.7)
dob date(2012-01-02, 2013-01-02),
age integer(4, 5),
name dict("first_names")
```

Field names are arbitrary, but field types must be particular values. See below for the complete list.
Some fields take arguments. Typically, a single numeric argument produces a field of that length. Whereas two values
signifies that you desire a random value in a particular range. The dict type requires the type of the
dictionary you are interested in as an argument.

#### Supported field types

* string
* decimal
* integer
* date
* dict (see [link](https://github.com/ThoughtWorksStudios/datagen/wiki/Dictionary-Field-Type) for more details)

#### Field type argument format/type

* string: integer
* decimal: float
* integer: (int, int)
* date (YYYY-MM-DD, YYYY-MM-DD)
* dict ("dictionary_type")


#### Inheriting from entities

```
Person: thing {}
```

Inherited entities will inherit all fields from their sub-entity, and will overwrite underlying fields with the same name.

#### Generating entities

```
generate (1, thing)
```

The generate keyword takes the number of entities and the entity name as arguments.

#### Overriding fields in generate statements

```
generate (5, Person { status "hmmm" })
```

You can pass comma-separated fields along with the entity name to override existing fields in a definition.

Note that if the field type is overriden, then the generate statement will generate the entity with a different name. for example:

```
generate (5, Person { status integer(1,10) })
```

Would output entities with a name like Person2356.

### Prerequisites

There are no prerequisites for running the binary.

### Building from source

The included Makefile has targets to get you started.

    make list
      build clean depend docker local release run test wercker


The simplest way to get started is to use docker. Install Docker for Mac and then run:

    make docker

This will create a docker container, build the software, and run the example file.

If you prefer to avoid containers, try:

    make local

## Running the tests

Simply run:

    make test
