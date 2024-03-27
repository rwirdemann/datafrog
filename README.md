# Database Dragon

## Development

Run `make` to create the two binaries `rt-record` and `rt-listen` in `$GOPATH/bin`.

```
$ make 
```

## Recording

Start recorder in a shell. Expects a `config.json` file in the current directory and writes matching
staments to the file passed as first command line argument.

```
$ rt-record create-job.rt
```

## Listening

Start listener in a shell. Expects a `config.json` file in the current directory and reads recorded
staments from the file passed as first command line argument.

```
$ rt-listen create-job.rt
```





