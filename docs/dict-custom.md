#Creating Custom Dictionaries

## Definition

Defining a basic custom dictionary is fairly straight forward.

1. Create a file named for the dictionary name you want. For example, if you want a dictionary of cat species (i.e. `dict("cats")`), you would create a file called cats (note that there is no extension). This will be where the list of possibilities will be located.
```bash
$ touch cats
```

2. At this point the cats dictionary is empty which isn't very useful. You can add possible values to the dictionary by adding the values to the dictionary file. Each value should be separated by a newline (\n). See the example bellow:
```bash
$ echo -e "lion\ndomestic\ntiger\npanther\nbobcat" > cats
$ cat cats
lion
domestic
tiger
panther
bobcat
```

At this point, you've successfully created a custom dictionary. You can now use this dictionary in your spec files.

```example-success
entity Person {
  favorite_animal: $dict("cats")
}

generate (2, Person)
```

Which should result in:

```bash
$ ./bobcat person.lang
$ cat entities.json
{
  "Person": [
    { "favorite_animal": "bobcat" },
    { "favorite_animal": "panther" }
  ]
}
```

**Note:** You can specify where bobcat will look for custom dictionary files using the `-d` argument. If it's not specified, then bobcat will look for dictionary files where the specification file is located.

## Dictionary Formats

Dictionary files with the filename suffix of `_format` contain a list of formats of possible values.

There are a few key components to dictionary formats:

### Numeric formats

The pound sign is used represent an arbitrary integer. So when you get a value from the dictionary each `#` will be replaced by in integer.

A simple example is the phone numbers dictionary:

```bash
$ cat phone_numbers_format
#-###-###-####
###-###-####
```

Now we can simply reference it by name:

```example-success
entity Person {
  phone: $dict("phone_numbers")
}

generate(3, Person)
```

Result:

```bash
$ ./bobcat person.lang
$ cat entities.json
{
  "Person": [
    {
      "phone": "1-123-432-1234"
    },
    {
      "phone": "321-654-9876"
    },
    {
      "phone": "7-666-666-1234"
    }
  ]
}
```

### Composite Formats

A composite format allows for the reference of other dictionaries. The possible values will be composed of values from referenced dictionaries. References are separated from other parts of the format by the reserved character `|`

A good example of a composite format is the full name dictionary ([full_names_format](https://github.com/ThoughtWorksStudios/bobcat/blob/master/dictionary/data/en/full_names_format)). It contains the following:

```bash
$ cat full_names_format
first_names| |last_names
```

Usage:

```example-success
entity Person {
  full_name: $dict("full_names")
}

generate (2, Person)
```

Result:

```bash
$ ./bobcat person.lang
$ cat entities.json
{
  "Person": [
    { "full_name": "Rick Sanchez" },
    { "full_name": "Morty Smith" }
  ]
}
```

#### Nesting Composite Formats
It's also possible to reference other composite format dictionaries as well:

```bash
$ cat thing_format
full_names_format| |c###
```

Usage:

```example-success
entity Person {
  thing: $dict("thing")
}

generate (2, Person)
```

Result:
```bash
$ ./bobcat person.lang
$ cat entities.json
{
  "Person": [
    { "thing": "Rick Sanchez c137" },
    { "thing": "Morty Smith c132" }
  ]
}
```
