# DataOwl (dbd)

A simple command line tool to record and validate database statements.

The current version was tested with MySQL on MacOS.

## Motivation

The basic idea of `dbd` is to record and replay usecase scenarios in order to
ensure that its repeated execution triggers the same database calls. One of the
challenges of this approach is dynamic data recognition like id's or date
fields. Consider the following example:

```
insert into job (description, ..., id) values ('World', ..., 1)
insert into job (description, ..., id) values ('World', ..., 2)
```

Both statements should be recognized as equal even if they slightly differ. To
detect and ignore this dynamic data we've build a token based diff algorithm
that considers a list of allowed differences when comparing two statemenst. This
algorithm requires to run the same usecase twice. The first run is called
`recording`and the second `verification`.

## Recording

```
dbd record --out create-job.rt
```

Starts the recorder in the command line. Hit enter and run the usecase scenario
through the UI of your application. Hit enter again when you are finsihed.   

## Verfication

```
dbd verfify --expectstions create-job.rt
```

Starts the verfier in the command line to verify the expectaitons recorded
during the previous run in `create-job.rt`. Hit enter and run exactly the same
usecase scenario again. The verification stops when all expectation have been
verfied successfully. It stops with an error when the verfication process got
out of order.

## Development

Run `make` to create the `dbd`binary.

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
rule thus only statments that contain `select job` but not `publish_trials<1`
are recorded.

## Further options

See help to see recording and validation options.

```
$ dbd --help
```