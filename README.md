# DataFrog (dfg)

A simple API to record and verify statement logged by external systems, e.g.
databases.

The current version was tested with MySQL on MacOS.

## Motivation

The basic idea of `dfg` is to record and replay usecase scenarios in order to
ensure that its repeated execution triggers the same database calls. One of the
challenges of this approach is dynamic data recognition like id's or date
fields. Consider the following example:

```
insert into job (description, ..., id) values ('World', ..., 1)
insert into job (description, ..., id) values ('World', ..., 2)
```

Both statements should be recognized as equal even if they slightly differ. To
detect and ignore this dynamic data we've build a token based diff algorithm
that considers a list of allowed differences when comparing two statements. This
algorithm requires to run the same usecase at least twice. The first run is
called `recording`. The second and all succeeding runs are called
`verification`.

## Development

Run `make` to create the `dfgapi` and `dfgweb` binaries.

```
$ make 
```

## Configuration

Expects a `config.json` file in the current or config subdirectory according to
the following format:

```json
{
  "sut": {
    "base_url": "http://localhost:8080"
  },
  "filename": "/usr/local/var/mysql/development.log",
  "logformat": "mysql",
  "patterns": [
    "insert into job",
    "update job",
    "delete",
    "select job!publish_trials<1"
  ],
  "expectations": {
    "report_additional": true
  },
  "ui_driver": "none",
  "web": {
    "port": 8081,
    "timeout": 120
  },
  "api": {
    "port": 3000
  }
}
```

A database statement is only recorded if it contains one of the configured
patterns. The pattern format `select job!publish_trials<1` contains an exclude
rule thus only statements that contain `select job` but not `publish_trials<1`
are recorded.

Allowed logformat: mysql | postgres

## API

Run `dfgapi` to start the backend.

```
# List of avaiable tests
GET /tests 

# Creates test 'name' and starts recording
POST /tests/{name}/recordings [POST]

# Stops recording of test 'name' 
DELETE /tests/{name}/recordings [DELETE]

# Delete test 'name'
DELETE /tests/{name}

# Starts verification of test 'name'
PUT /tests/{name}/verifications [PUT]

# Stops verification of test 'name'
DELETE /tests/{name}/verifications 
```

## Web UI

Run `dfgweb` to start the web frontend. Requires a running backend.