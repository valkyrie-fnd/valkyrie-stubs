# Valkyrie stubs

Stubs for valkyrie integration testing. Integration tests for |Valkyrie](https://github.com/valkyrie-fnd/valkyrie) use this library to cut out the dependency to a PAM. 

## Running standalone
Another option is to start both `Valkyrie` and `Valkyrie-stubs` in standalone mode in order to 
test some provider integration. 

To start standalone instances of Valkyrie with stubs using a memory database:
```bash
git clone git@github.com:valkyrie-fnd/valkyrie-stubs.git
cd valkyrie-stubs
docker-compose up 
```

## Running tests with fault injection
When starting the stubs in standalone mode a web interface is available at http://localhost:8080/broken. 

There some predefined error scenarios can be triggered, like connection issues.  

In order to extend the available cases take a look at [scenarios.go](./broken/scenario.go).

## A few curls
```bash
# Create a session using the backdoor
SID=$(curl -s -H 'Content-type:application/json' -d '{"sid":"sid2", "userId":"5000001"}' 'localhost:3000/backdoors/evolution/sid?authToken=evo-api-key' | jq -r '.sid')

# Get balance
curl -i -H "Authorization:Bearer pam-api-token" \
    -H "X-Player-Token:$SID" -H "X-Correlation-ID:sfasdfa" \
    -X GET "localhost:8080/players/5000001/balance?provider=evolution&currency=EUR"

# Place bet
curl -i "localhost:8080/players/5000001/transactions?provider=evolution" -H "Authorization:Bearer pam-api-token" \
    -H "X-Player-Token:$SID" -H "X-Correlation-ID:sfasdfa" -H 'Content-type:application/json' \
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
