#! /bin/bash
mongoimport -h localhost:27017 --drop -d quotations --jsonArray --file bobcat/output/flat_quotes.json
mongoimport -h localhost:27017 --drop -d quotations --jsonArray --file bobcat/output/users.json
mongoimport -h localhost:27017 --drop -d quotations --jsonArray --file bobcat/output/nested_quotes.json
