When defining a field one can specify a "count range" to indicate that a field should produce an array of 0 or more values. The count range syntax is a range (lower-bound-number, followed by `..`, followed by upper-bound-number), surrounded by angled brackets (`<`, `>`).

```
# the `emails` field will yield an array of 1 - 5 email addresses
(the number of emails in the array will vary from one JSON object to another)

entity {
  emails: $dict("email_address")<1..5>
}
```
To produce an array with an exact number of values, set the lower bound and upper bound to the same value.

```
# the `login` field will yield an array of 2 logins

entity {
  logins: $dict("logins")<2..2>
}
```
Count ranges can be used with built-in field types (excluding distribution), but they can also be used with entity field types. This can be particularly useful when creating one to many relationships.

```
# the `items` field will yield an array of 0 to 10 CartItems (or IDs if --flatten output)

entity {
  items: CartItem<0..10>
}
```

Note count ranges cannot be used in conjunction with the [unique](builtins.md#unique-value-flag) flag.
