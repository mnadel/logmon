# Logmon

A logfile monitor.

# Description

Logmon will scan your logfiles for errors, collect all new errors since the last time it was ran, and then send those errors somwhere (e.g. Syslog, SMTP, REST, etc).

# Configuration

Create a JSON doc that's parseable into:

    type Configuration struct {
        AlertConfig map[string]string `json:"alert"`
        Db          string            `json:"db"`
        Logs        []string          `json:"logs"`
        ErrorTokens []string          `json:"toks"`
    }

`AlertConfig` is a map of config params which are passed to the alerter.

`Logs` is a list of file globs.

`Db` is the path to a BoltDB database file; it'll be created the first time Logmon is ran.

`ErrorTokens` is a list of tokens, any of which identify a line in your logfile as an error; defaults to `[ERROR,FATAL]`.

# SMTP Alerter

If the `alert` config object looks like the below, an email will be sent out with a list of errors.

    {
        "smtp": "smtpserver:port",
        "from": "From email address",
        "to": "To email address",
        "subject": "Subject of email"
    }

# Dependencies

1. [BoltDB](https://github.com/boltdb/bolt) `go get github.com/boltdb/bolt`
1. [imohash](https://github.com/kalafut/imohash) `go get github.com/kalafut/imohash`

# Implementation notes

We use `imohash` to create a md5-like hash of each logfile in constant time. This allows us to quickly and reliably detect if a file changed since the last time we ran.

We store both the hash and the last point in the file that we read in BoltDB.

On each run, if we detect that a file changed, we'll seek to the last byte read, and will scan for new errors.

# TODO

1. Additional/pluggable alerters
