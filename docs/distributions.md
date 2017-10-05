## Distributions

Distributions allow you the specify the shape that the generated data should take. Currently, there are a few distributions that are builtin to bobcat. In the future we intend to have a way for users to define their own distributions.

**NOTE: Distributions may only be used in entity field declarations.**

The following are currently supported, builtin distributions:

| Name       | Description                                     | Usage                                                        |
|------------|-------------------------------------------------|--------------------------------------------------------------|
| `*normal`  | The [normal gaussian distribution](https://en.wikipedia.org/wiki/Normal_distribution) | Takes 2 numeric values as boundaries, e.g. `[1, 100]`. Outputs floating point values. |
| `*weight`  | Distributes values based on probability weights | [[`weightedArg1`](#weighted-arguments), ..., `weightedArgN`] |
| `*percent` | Similar to `*weight` where the weights must sum to `1.0` | [[`weightedArg1`](#weighted-arguments), ..., `weightedArgN`] |

### Uniform Distribution

Bobcat's [native functions](../README.md#native-functions) already follow a uniform distribution (i.e. equal probability for any value in the domain).

### Weighted Arguments

A weighted argument is in the following format: `left-hand-expression => right-hand-expression`

The **left hand expression** must evaluate to a numeric value, and will be evaluated immediately during entity declaration. This implies that one cannot reference another field in the entity when specifying weights (this would not be meaningful anyway).

The **right hand expression** may be any expression, and will evaluate every time an entity's field value is generated. This implies that the right hand expression may reference other fields within the entity. The right hand expressions need not return the same type for each weighted argument within a distribution.

### Example

```example-success
let MIN_WEIGHT = 5
let MAX_WEIGHT = 400

let BASE = 50

entity User {
  name: $dict("full_names"),

  age: *percent ~ [
    0.25 => $int(1, 18),
    0.5 => $int(18, 50),
    0.25 => $int(50, 80)
  ],

  favorite_flavor: *weight ~ [
    BASE => "vanilla",
    10 * BASE => "chocolate",
    BASE / 5 => "strawberry",
    BASE / 25 => null,
  ],

  weight: *normal ~ [MIN_WEIGHT, MAX_WEIGHT],
}
```
