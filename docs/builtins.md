### Enumerated Field (`$enum`)

Enumerated values are sort of like inlined $dictionaries. `$enum(collection)` picks a value from the given collection:

```example-success
# declare a collection
let statuses = ["To do", "Doing", "Done!"]

entity Work {
  status: $enum(statuses) # randomly picks from statuses
}
```

`generate()` statements also yield collections of `$id`s from generated entities. This can be used in conjunction with `$enum` fields to define relationships:

```example-success
entity CatalogItem {
  name: $str(),
  sku:  $int(1000, 3000)
}

# Assign the collection from generate() to a variable
let Catalog = generate(20, CatalogItem)

# each cart will have 1 - 5 CatalogItems as its contents
entity ShoppingCart {
  contents: $enum(Catalog)<1..5>
}

```

