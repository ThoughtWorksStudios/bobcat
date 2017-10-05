### Primary Key Statement
When an enitity is [generated](../README.md#generating-entities-generate-expressions) an $id field (type $uid) is automatically included in the resulting JSON object. This field is configurable via a primary key statement `pk(<field-name>, <field-type>)`. The `field-name` can be any identifier and the `field-type` can be either $uid or $incr.

```example-success
pk("Id", $incr)
```

Primary key statements defined within an entity should be the first statement after the curly brace before the suceeding fields. The primary key will only be changed for that entity or any entities extending that entity.

```example-success
entity Product {
  pk("productId", $incr)  # will replace $id with "productId" for Products only
  name: $dict("nouns"),
  quantity: $int()
}
```

Primary key statements outside of an entity will apply to all subsequent entity definitions.

```example-success
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

