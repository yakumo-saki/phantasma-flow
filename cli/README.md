# phantasma-flow-cli

Work in progress... (maybe WIP forever)

## What is phantasma flow?

* See https://github.com/yakumo-saki/phantasma-flow

# usage

`phctl command subcommand`

## exit code

| num | reason  |
| ------- | -------- |
| 0       | Success  |
| 4 to 15 | There is some warning |
| 16 and over | Failure |
# command and subcommands

## command: ping

`phctl ping [auth]`

Tests phantasma-flow server connectivity.
This not tests authenticate.

### exit codes

* 0 success
* 16 failure

### subcommand

#### auth

Also check authentication.

## command: get

`phctl get [jobs|nodes]`

### exit codes

* 0 success
* 16 failure

### General options

* `-o [text|json]` specify output type.

### Subcommands

#### jobs

Get all jobs

#### nodes

## command: log

`phctl log [-f] JobId or RunId`

### exit codes

* 0 success
* 16 failure

### General options

* `-f` follow log outputs.

### Limitations

* Log output may lack some lines, when viewing running job log.
