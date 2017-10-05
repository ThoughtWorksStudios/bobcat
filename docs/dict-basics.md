# Dictionaries

The dictionary field type can be used to generate a value from a list of possibilities. This tool offers a handful of builtin dictionaries, and is able to use custom dictionaries that you can define yourself.

### Builtin dictionary types

The following are the out of the box dictionaries:

* adjectives
* characters
* cities
* colors
* companies
* continents
* countries
* currencies
* currency_codes
* domain_zones
* email_address
* first_names
* full_names
* genders (Note: the gender dictionary is non-binary inclusive ^^)
* industries
* jobs
* languages
* last_names
* months
* months_short
* name_prefixes
* name_suffixes
* nouns
* patronymics
* phone_numbers
* state_abbrevs
* states
* street_address
* street_suffixes
* streets
* streets
* top_level_domains
* weekdays
* weekdays_short
* words
* zip_codes

Example:
```example-success
entity Person {
  name: $dict("first_names")
}

generate(2, Person)
```

Result:
```bash
$ cat entities.json
{
  "Person": [
    {
      "name": "Jet"
    },
    {
      "name": "Alex"
    }
  ]
}
```
