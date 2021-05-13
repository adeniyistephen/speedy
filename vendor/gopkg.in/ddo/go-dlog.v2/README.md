# go-dlog [![Doc][godoc-img]][godoc-url]
:pager: conditional logging for Golang libraries &amp; applications

![screenshot][screenshot]

[godoc-img]: https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat-square
[godoc-url]: https://godoc.org/gopkg.in/ddo/go-dlog.v2
[screenshot]: http://i.imgur.com/cZOEREE.png

## installation

```sh
go get gopkg.in/ddo/go-dlog.v2
```

## usage

```go
logger := dlog.New("logger", nil)

logger.Debug("some log")
logger.Info("some log")
logger.Warn("some log")
logger.Error("some log")
```

```sh
DLOG=* go run example.go
```

## env

set ``DLOG`` environment(case-insensitive) to

* ``DEBUG`` or ``*`` to enable ``#Debug`` and above
* ``INFO`` to enable ``#Info``, ``#Done``, ``#Fail`` and above
* ``WARN`` to enable ``#Warn`` and above
* ``ERROR`` to enable ``#Error``
