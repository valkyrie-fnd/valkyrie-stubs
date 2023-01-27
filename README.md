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
