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

There is also a `docker-compose.yml` available in the repository root which will start **valkyrie** together with **valkyrie-stubs** 
acting as PAM:

```bash
docker-compose up 
```

Valkyrie stubs may also be used as a testing library, for example by starting a genericpam server programmatically as
in this Valkyrie test suite [example](https://github.com/valkyrie-fnd/valkyrie/blob/main/provider/internal/test/suite.go#L51).

### Test configuration
The stubs contain an in-memory database. To configure the database, edit the datastore [config file](./datastore.config.yaml). Make sure that the appropriate player, provider, provider API key and game are present in the config. Otherwise add the values.

To configure the valkyre module for the test, make sure that the provider is present [here](./valkyrie_config.yml) too.

*Note*: The config files mentioned in this section are included in the image in the docker build, so there is no point in changing them in runtime. If stubs is built and run as a go binary, the config files should be present in the execution directory and can be updated between each stop and start.

### Running tests with fault injection using Broken

When starting the stubs in standalone mode a web interface is available at http://localhost:8080/broken. 

There some predefined error scenarios can be triggered, like connection issues.  

In order to extend the available cases take a look at [scenarios.go](./broken/scenario.go).

Once started, use the web interface to inject an appropriate error. If repeated errors (of the same kind) are wanted, start as many tabs of the web interface as wanted and inject the fault once per tab.

Each fault will be triggered and subsequently reset by each wallet request (i.e. balance or transaction). To make sure the errors work as intended, the curl calls below might become handy. For proper game tests, run stubs together with valkyrie as described above, obtain a session token (see curl below) and fire wallet requests towards the appropriate provider endpoints.

Expected outcome of the tree boiler plate scenarios provided for balance and withdrawal requests are:
* Requests directly to stubs PAM
  * *Undefined Error* - 500 Internal Server Error. Body: {"error":{"code":"PAM_ERR_UNDEFINED","message":"forced error"},"status":"ERROR"}
  * *Timeout 5s* - 408 Request Timeout. Body: N/A
  * *Connection close* - No response
* Requests to provider endpoints (Evolution in this example)
  * *Undefined Error* - 200 OK. Body: {"status":"UNKNOWN_ERROR","balance":0.000000,"bonus":0.000000,"uuid":"123"}
  * *Timeout 5s* - 200 OK. Body: {"status":"UNKNOWN_ERROR","balance":0.000000,"bonus":0.000000,"uuid":"123"}
  * *Connection close* - 200 OK. Body: {"status":"UNKNOWN_ERROR","balance":0.000000,"bonus":0.000000,"uuid":"123"}

*Note*: To trigger the connection close error, two instances of broken need to be activated. Valkyrie provider software performs one retry in case of connection close error.

### A few curls

*Note*: The curl commands in these examples are executed directly towards the PAM. To test a provider game together with Valkyrie and valkyrie-stubs, provider specific endpoints are used.

Create a session for test using the backdoor HTTP API (note that this session token can be used in provider tests as well):

```bash
# Evolution specific
SID=$(curl -s -H 'Content-type:application/json' -d '{"sid":"sid2", "userId":"5000001"}' 'localhost:3000/backdoors/evolution/sid?authToken=evo-api-key' | jq -r '.sid')

# General
SID=$(curl -s -H 'Content-type:application/json' -d '{"provider":"caleta", "userId":"5000001"}' 'localhost:3000/backdoors/session' | jq -r '.result.token')

```

Get the balance using the first session (in this example) from the genericpam stub:

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
