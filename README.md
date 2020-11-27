[![Go Report Card](https://goreportcard.com/badge/github.com/cimomo/portfolio-go)](https://goreportcard.com/report/github.com/cimomo/portfolio-go)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/cimomo/portfolio-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/cimomo/portfolio-go/blob/master/LICENSE)

# portfolio-go: Portfolio Management for Geeks

A terminal based portfolio tracking, analysis and optimization tool implemented in Go. One screenshot is worth a thousand words:

![Screenshot](./examples/screenshots/strategic.png "Portfolio-go screenshot")

## Using portfolio-go

Start the program by running:
```
portfolio --profile <path-to-profile>
```

Here is a sample profile:
```
cash:
  value: 10000.00
  allocation: 10
portfolios:
- portfolio: FAAMG
  allocation: 90
  holdings:
  - symbol: FB
    quantity: 715
    allocation: 20
    basis: 20000
    watch: 230
  - symbol: AAPL
    quantity: 1172
    allocation: 20
    basis: 20000
    watch: 100
  - symbol: AMZN
    quantity: 78
    allocation: 20
    basis: 20000
    watch: 2800
  - symbol: MSFT
    quantity: 861
    allocation: 20
    basis: 20000
    watch: 200
  - symbol: GOOG
    quantity: 56
    allocation: 20
    basis: 20000
    watch: 1400
```
More examples can be found [here](examples/).