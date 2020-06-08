# FreeFeed database anonymizer

This tool anonymizes the [FreeFeed](https://github.com/FreeFeed/freefeed-server)
database dump.

It takes the `pg_dump` result (in plain SQL format) and anonymizes it in the 
following way:

 * Some (configurable) tables are keeping untouched.
 * Some (configurable) tables are cleaning up and becoming empty.
 * In all other tables:
    * All UUIDs are irreversibly encrypting. They remain correct
    UUIDs, however it is impossible to restore the original UUIDs back.
    * Some (configurable) column values are replacing by lorem ipsum-like content.
    This applies to texts, names, emails etc. There a several rules for different 
    types of content.
    * Some (configurable) column values are replacing by the predefined constant
    values.
 
## Build

Use [Go 1.14+](https://golang.org/) to build this program:
```
go get github.com/davidmz/anonymize-db
```

## Use

```
Flags of anonymize-db:
  -i, --input string    input file name (STDIN by default)
  -o, --output string   output file name (STDOUT by default)
  -r, --rules string    rules file name (JSON)
```

The only required option is the path to [rules.json](./rules.json) file.
This file contains specific anonymizing rules for different database tables.

The simplest usage is:

```
pg_dump -d freefeed -U freefeed -F p | anonymize-db -r rules.json > anonymized.sql
```

## Configure

The [rules.json](./rules.json) file contains specific anonymizing rules for 
different database tables. The file contains a mapping between the full-qualified 
table names (including schema names, like "public.users") and instructions.

The `{ "action": "KEEP" }` instruction means that the table must be keeping untouched.

The `{ "action": "CLEAN" }` instruction means that the table must be truncated and 
become empty.

The `{ "columns": {...} }` instruction deifies per-column anonymizing rules. Any column
that are not listed here (and any tables that are not listed in the file) will be
processed in general way: only the UUIDs will be encrypted.

The per-column rules are:

* `text` — long texts like posts and comments bodies;
* `shorttext` — short texts like names and titles;
* `uniqword` — unique words like usernames (the same username always converted
  to the same word);
* `uniqemail` — unique email addresses;
* `set:VALUE` — constant value to write to the column.

The `set:` value format is the same as the 
[PostgreSQL COPY format](https://www.postgresql.org/docs/current/sql-copy.html#id-1.9.3.55.9.2),
so to set a NULL value use `set:\N` rule (and add the second backslash for JSON). 

Only non-NULL and non-empty columns processed.

### Provided configuration

The provided [rules.json](./rules.json) is compatible with FreeFeed DB scheme with 
the last migration `20200422122956_multi_homefeeds.js`.

The `hashed_password` value in the `public.users` table is a hashed 'tester' string 
(so all users in resulted database will have the same 'tester' password).

## Debug

Set `DEBUG=anon` environment variable to see some debug information during the process.