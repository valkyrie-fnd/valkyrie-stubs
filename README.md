# Valkyrie Stubs

[![](https://img.shields.io/badge/License-MIT%20-brightgreen.svg)](./LICENSE.md)
[![](https://img.shields.io/github/actions/workflow/status/valkyrie-fnd/valkyrie-stubs/gh-workflow.yml)](https://github.com/valkyrie-fnd/valkyrie-stubs/actions/workflows/gh-workflow.yml)
![](https://img.shields.io/github/last-commit/valkyrie-fnd/valkyrie-stubs)
[![](https://img.shields.io/website?url=https%3A%2F%2Fvalkyrie.bet)](https://valkyrie.bet/docs)
![](https://img.shields.io/github/go-mod/go-version/valkyrie-fnd/valkyrie-stubs)
![](https://img.shields.io/github/languages/top/valkyrie-fnd/valkyrie-stubs)
![](https://img.shields.io/tokei/lines/github/valkyrie-fnd/valkyrie-stubs)
![](https://img.shields.io/maintenance/yes/2023)

## A testing library for the Valkyrie iGaming Aggregator

This repository contains testing stubs and utilities for [Valkyrie](https://github.com/valkyrie-fnd/valkyrie), an 
open source iGaming Aggregator.

This testing repository includes:

* **genericpam** - an implementation of Valkyrie's Player Account Management (PAM) API, for more information please read
  the [documentation](https://valkyrie.bet/docs/wallet/valkyrie-pam/valkyrie-pam-api)
* **memorydatastore** - a datastore implementation for **genericpam** persisted in memory, useful for building integration 
  tests
* **backdoors** - HTTP APIs for creating sessions and other testing data fixtures in genericpam
* **broken** - an HTTP API and a web interface for simulating various failure scenarios (timeouts, errors, etc)

## Building

The project is built from source by running:

```bash
go build
```

## Running

The project can run standalone using the built binary:

```bash
./valkyrie-stubs
```

There is also a `docker-compose.yml` available which will start **valkyrie** together with **valkyrie-stubs** 
acting as PAM:

```bash
docker-compose up 
```

Valkyrie stubs may also be used as a testing library, for example by starting a genericpam server programmatically as
in this Valkyrie test suite [example](https://github.com/valkyrie-fnd/valkyrie/blob/main/provider/internal/test/suite.go#L51).

### Running tests with fault injection using Broken

When starting the stubs in standalone mode a web interface is available at http://localhost:8080/broken. 

There some predefined error scenarios can be triggered, like connection issues.  

In order to extend the available cases take a look at [scenarios.go](./broken/scenario.go).

### A few curls

Create a session for the Evolution provider using the backdoor HTTP API:

```bash
SID=$(curl -s -H 'Content-type:application/json' -d '{"sid":"sid2", "userId":"5000001"}' 'localhost:3000/backdoors/evolution/sid?authToken=evo-api-key' | jq -r '.sid')
```

Get the balance using this session from the genericpam stub:

```shell
curl -i -H "Authorization:Bearer pam-api-token" \
    -H "X-Player-Token:$SID" -H "X-Correlation-ID: some-uuid-identifier1" \
    -X GET "localhost:8080/players/5000001/balance?provider=evolution&currency=EUR"
```

Place a bet (withdrawal transaction) using the same session:

```shell
curl -i "localhost:8080/players/5000001/transactions?provider=evolution" -H "Authorization:Bearer pam-api-token" \
    -H "X-Player-Token:$SID" -H "X-Correlation-ID: some-uuid-identifier2" -H 'Content-type:application/json' \
    -X POST -d '
    {
        "currency": "EUR",
        "provider": "evolution",
        "transactionType": "WITHDRAW",
        "cashAmount": 5,
        "providerTransactionId": "187",
        "transactionDateTime": "2022-02-25T10:11:43.511Z",
        "providerGameId": "vctlz20yfnmp1ylr",
        "providerRoundId": "ABC001"
    }'
```
