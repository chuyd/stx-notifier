# A simple gerrit notifier

This simple program queries the Openstack REST API to get the latest change from
the `stx-tools` repository and displays a notification in the notification bar.

## Building

To build first download the `logrus` dependency.

```
go get github.com/sirupsen/logrus
```

then just run:

```
go build
```

## Running

Just launch the program.

```
./stx-notifier
```

This will loop forever checking for new changes every minute.

## TODO

- Support more projects
