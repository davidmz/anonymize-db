# FreeFeed database anonymizer

This tool anonymizes the [FreeFeed](https://github.com/FreeFeed/freefeed-server)
database dump.

It takes the `pg_dump` result (in plain SQL format) and anonymizes it in the 
following way:

 * Some (configurable) tables are keeping untouched.
 * Some (configurable) tables are cleaning up and becoming empty.
 * In all other tables:
    * Some or all (configurable) UUIDs are irreversibly encrypting. They remain correct
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
  -c, --config string   config file name (JSON)
  -i, --input string    input file name (STDIN by default)
  -o, --output string   output file name (STDOUT by default)
```

The only required option is the path to config file. This file contains specific anonymizing rules
for different database tables. There are two predefined configuration files in this repository:

  * [config.anon.json](./config.anon.json) — fully anonymizes database with all UUIDs and user data;
  * [config.candy.json](./config.candy.json) — anonymizes database for the stage server
    (candy.freefeed.net): usernames, hashed passwords and UUIDs are not changes (so users 
    can sign in with their freefeed credentials), but all texts and attachments are anonymizes;

The simplest usage is:

```
pg_dump -d freefeed -U freefeed -F p | anonymize-db -c config.anon.json > anonymized.sql
```

## Configure

The config file contains specific anonymizing rules for different database tables. The file contains
a mapping between the full-qualified table names (including schema names, like "public.users") and
instructions.

The top-level `encryptUUIDs` boolean value defines whether the all database UUIDs should be encrypted.

The `{ "keep": true }` per-table instruction means that the table must be keeping untouched.

The `{ "clean": true }` per-table instruction means that the table must be truncated and 
become empty.

The `{ "columns": {...} }` per-table instruction defines the per-column anonymizing rules. Any column
that are not listed here (and any tables that are not listed in the file) will be
processed in general way: only the UUIDs will be encrypted if the top-level `encryptUUIDs` is true.

The per-column rules are:

* `text` — long texts like posts and comments bodies;
* `shorttext` — short texts like names and titles;
* `uniqword` — unique words like usernames (the same username always converted
  to the same word);
* `uniqemail` — unique email addresses;
* `uuids` — encrypt any UUIDs in column (if the global `encryptUUIDs` parameter isn't set);
* `set:VALUE` — constant value to write to the column.

The `set:` value format is the same as the 
[PostgreSQL COPY format](https://www.postgresql.org/docs/current/sql-copy.html#id-1.9.3.55.9.2),
so to set a NULL value use `set:\N` rule (and add the second backslash for JSON). 

Only non-NULL and non-empty columns processed.

### Provided configurations

The provided configs are compatible with FreeFeed DB scheme with 
the last migration `20210224194435_comment_numbers.js`.

The `hashed_password` value in the `public.users` table ([config.anon.json](./config.anon.json)) is
a hashed 'tester' string (so all users in resulted database will have the same 'tester' password).

Both configs are cleans up the posts and comments `body_tsvector` fields. You must rebuild the
search index after the anonymization if you want to use search.


## Debug

Set `DEBUG=anon` environment variable to see some debug information during the process.