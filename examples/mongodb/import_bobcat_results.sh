#! /bin/bash
(
  cd bobcat
  mkdir -p output
  ../../../bobcat-darwin -o output/nested_quotes.json nested_quote.lang
  ../../../bobcat-darwin -o output/e.json -s flat_quote.lang
  mv output/e-user.json output/users.json
  mv output/e-quotes.json output/flat_quotes.json
)

mongoimport -h localhost:27017 --drop -d quotations --jsonArray --file bobcat/output/flat_quotes.json
mongoimport -h localhost:27017 --drop -d quotations --jsonArray --file bobcat/output/users.json
mongoimport -h localhost:27017 --drop -d quotations --jsonArray --file bobcat/output/nested_quotes.json
