# DataGen
[![wercker status](https://app.wercker.com/status/98be3d80966b1a3a006c0465c76aa8ef/s/master "wercker status")](https://app.wercker.com/project/byKey/98be3d80966b1a3a006c0465c76aa8ef)

A data generation tool. Just define concepts in our input file format, and the tool will generate JSON objects that can be used to insert realistic-looking data into your application.

## Getting Started

### Quickstart

1. Download the latest [release](https://github.com/ThoughtWorksStudios/datagen/releases)
2. Run the tool over the sample input file:

        ./datagen example.lang > my_data.json

3. Modify the sample file or create one from scratch to generate your own custom entities

### Input file format

```
def Person {
  full_name dict("full_name"),
  login string(4),
  thing string,
  dob date(1985-01-02, 2000-01-01),
  age  decimal(4.2, 42.7)
}

def Cat {
  paws string(2),
  dob date(2012-01-02, 2013-01-02),
  age integer(4, 5),
  lives integer(0, 9),
  name dict("first_name")
}

generate Person(50)
generate Cat(25)
```

The input file contains definitions of entities (the objects, or concepts found in your software system), fields on those entities (properties that an entity posses), and a 'generate'
keyword to produce the desired number of entities in the resulting JSON output. An entity has an arbitrary name,
as do fields. The only other concept in this system is that of a dictionary, which is used to provide
realistic values for fields that would otherwise be difficult to generate data for (like a person's name).

#### Defining entities

```
def Person
```

The 'def' keyword is required, but the name after the def can be one of your choosing. There are no predefined entity types.

#### Defining fields

```
login string(4),
age  decimal(4.2, 42.7)
dob date(2012-01-02, 2013-01-02),
age integer(4, 5),
name dict("first_name")
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
* dict

#### Field type argument format/type

* string: integer
* decimal: float
* integer: (int, int)
* date (YYYY-MM-DD, YYYY-MM-DD)
* dict ("dictionary_type")

#### Supported dictionary types

The following is a list of supported dictionary types:

* last_name
* first_name
* city
* country
* state
* street
* address
* email
* zip_code
* full_name
* random_string

#### Generating entities

```
generate Person(50)
```

An entity generator only supports a single argument at this time, which is the number of entities that you'd like the program to produce.
### Prerequisites

There are no prerequisites for running the binary, but if you want to build the code, you'll need the [latest Go runtime](https://golang.org/dl/).

### Building from source

First, install the Go language on your target platform and add $GOPATH to your $PATH. The default location for $GOPATH is ~/go/bin. Then run the default target in the Makefile.

```
make
```

This will produce a binary called 'datagen'. If you've added $GOPATH to your $PATH, this binary is available from anywhere on the filesystem.

## Running the tests

Simply run:

        go test ./...
