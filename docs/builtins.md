### Customizing date formats

The built-in date field (i.e. `$date(min, max, format)`) can take an optional 3rd argument: a `strftime` format $string, e.g. `"%b %d, %Y %H:%M:%S"` All dates generated for that field will adhere to the format provided.

### Enumerated Field (`$enum`)

Enumerated values are sort of like inlined $dictionaries. `$enum(collection)` picks a value from the given collection:

```
# declare a collection

let statuses = ["To do", "Doing", "Done!"]

entity Work {
  status $enum(statuses) # randomly picks from statuses
}
```

`generate()` statements also yield collections of `$id`s from generated entities. This can be used in conjunction with `$enum` fields to define relationships:

```
entity CatalogItem {
  name: $string,
  sku:  $int(1000, 3000)
}

# Assign the collection from generate() to a variable
let Catalog = generate(20, CatalogItem)

# each cart will have 1 - 5 CatalogItems as its contents
entity ShoppingCart {
  contents: $enum(Catalog)<1..5>
}

```

### Distribution Field

Distribution fields allow you the specify the shape that the generated data should take. Currently, there are a few supported distributions that are builtin to bobcat. In the future we intend to have a way for users to define their own distributions.

The following are currently supported, builtin distributions:

| Name      | Value                                                                                | Allowed Fields | Format                                      |
|-----------|--------------------------------------------------------------------------------------|----------------|---------------------------------------------|
| `"normal"`  | The [normal gaussian distribution](https://en.wikipedia.org/wiki/Normal_distribution)| $float         | ("normal", $float(..), $float(..), ...)       |
| `"uniform"` | A uniform distribution                                                               | $int, $float   | ("uniform", $int(..), $int(..), ...)          |
| `"percent"` | specify the % something should occur                                                 | all            | ("percent", x% => field(..), y% => field(..)) |
| `"weight"`  | probability weights                                                                  | all            | ("weight", x => field(..), y => field(..))    |


example:
```
entity User {
  name: $dict("full_names"),
  age: $distribution("percent",
    25% => $float(1.0, 18.0),
    50% => $float(18.0, 50.0),
    25% => $float(50.0, 80.0)
  ),
  favorite_number: $distribution("weight",
    55  => $int(1, 15),
    500 => $int(15, 30),
    2   => $int(30, 80)
  ),
  weight: $distribution("normal", $float(1.0, 400.0)),
  status: $distribution("percent",
    10% => $enum(["disabled"]),
    90% => $enum(["pending", "active"])
  ),
  email: $dict("email_address"),
  email_confirmed: $distribution("percent",
    50% => "yes",
    50% => "no"
  )
}
```

### Primary Key Statement
When an enitity is [generated](../README.md#generating-entities-generate-expressions) an $id field (type $uid) is automatically included in the resulting JSON object. This field is configurable via a primary key statement `pk(<field-name>, <field-type>)`. The `field-name` can be any identifier and the `field-type` can be either $uid or $incr.

```
pk("Id", $incr)
```

Primary key statements defined within an entity should be the first statement after the curly brace before the suceeding fields. The primary key will only be changed for that entity or any entities extending that entity.

```
entity Product {
  pk("productId", $incr)  # will replace $id with "productId" for Products only
  name: $dict("nouns"),
  quantity: $int()
}
```

Primary key statements outside of an entity will apply to all subsequent entity definitions.

```
entity User {  # will have an $id field
  login: $dict("email_address"),
  status: $enum(["enabled", "disabled", "pending"])
}

pk("id", $incr)

entity Admin << User { # will have an $id field because extends previously defined User
  superuser: true
}

entity Book { # will have an "id" field due to pk statement
  title: $dict("nouns"),
  author: $dict("full_names")
}
```

