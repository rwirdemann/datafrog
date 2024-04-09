# Database Dragon (dbd)

A simple command line tool to record and validate database statements.

The current version was tested with MySQL on MacOS.

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

## Recording

The basic idea of `dbd` is to record and replay usecase scenarios in order to
ensure that its repeated execution triggers the same database calls. One of the
challenges of this approach is dynamic data recognition like id's or date
fields. Consider the follwing example:

```
insert into job (description, ..., id) values ('World', ..., 1)
insert into job (description, ..., id) values ('World', ..., 2)
```

Both statements should be recognized as equal even if they slightly differ. To
detect and ignore this dynamic data we've build a token based diff algorithm
that considers a list of allowed differences when comparing two statemens. This
algorithm requires to run the same usecase twice during the recording phase of
the test:

```
# 1. start recording
dbd record --out create-job.rt

# 2. manually run the usecase scenario through its ui

# 3. start the verfication
dbd verify --expectations create-job.rt

# 4. manually run the usecase scenario a second time
```

## Further options

See help to see recording and validation options.

```
$ dbd --help
```