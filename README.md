# DataFrog (dfg)

A simple command line tool to record and verify statement logged by external
systems, e.g. databases.

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

## Recording

```
dfg record --out create-job.json
```

Starts the recorder in the command line.  

## Verification

```
dfg verify --expectations create-job.json
```

Starts the verifier in the command line to verify the expectations recorded
during the previous run in `create-job.json`. The expectations file is updated
according to the fulfilled expectations.

## Development

Run `make` to create the `dfg` binary.

```
$ make 
```

## Configuration

Expects a `config.json` file in the current directory according to the following
format:

```json
{
  "filename": "/usr/local/var/mysql/development.log",
  "patterns": [
    "insert into job",
    "update job",
    "delete",
    "select job!publish_trials<1"
  ]
}
```

A database statement is only recorded if it contains one the configured
patterns. The pattern format `select job!publish_trials<1` contains an exclude
rule thus only statements that contain `select job` but not `publish_trials<1`
are recorded.

## Further options

See help to see recording and verfication options.

```
$ dfg --help
```