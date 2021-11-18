# Demo Bank App
This is a demo http application written in Golang

## Architecture
The app uses the layered architecture: api, service and storage.
The project structure has the source code and tests located in separate folder: `src` and `test` correspondingly.

### Database access/database migrations
The database is Postgresql.
The schema is set up using the database migration library [sql-migrate](https://github.com/rubenv/sql-migrate).
These migrations run on start up, so the schema is up-to-date right after the service is running.

### Authentication/authorization
The authentication is implemented using a mock service. 
This mock assumes that there are 2 valid users:
- user with id `1` with token `token_user_1`
- user with id `2` with token `token_user_2`

All endpoints are covered with the authentication using a Bearer token, so the requests have to contain the header
`Authorization: Bearer ...`.
Requests with no token have `401` http code, and with unknown tokens - `403`.

The authorization compares the account owner id with the user from the token.
If the user is not an owner, the request is rejected with `403`.

### Error handling
The business errors are defined in the folder `./src/errors` for each specific corner case situation.
They are pattern matched in the `./src/api/account.go`.
Each error will be serialized using the same format and specific http code. 
For example, requesting a not existing account would result in `404` and the response body
```json
{"message":"The account 6 does not exist"}
```

## Running tests
The prerequisite for running tests is golang 1.17 and docker.
Docker is required, because the storage tests call the real Postgresql.

This command starts docker, runs tests and then stops docker
```shell
make run_test
```

## Building the docker image
The docker image is build using the multi-stage docker build.
```shell
make build_docker
```

## Running the demo locally

### Starting the server
Assuming that the docker image was built in the previous step,
the demo can be started using docker compose.

The application fails fast and restarts, until Postgresql is reachable.

```shell
docker compose up
```

### Sample cUrl requests

1) Create an account for the user 1
```shell
curl --header 'Authorization: Bearer token_user_1' --request POST 'http://localhost:8000/accounts'
```
2) Check the balance for the first account
```shell
curl --header 'Authorization: Bearer token_user_1' 'http://localhost:8000/accounts/1'
```
3) Top up the first account
```shell
curl --request POST 'http://localhost:8000/top-up' \
--header 'Authorization: Bearer token_user_1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": 1,
    "amount": 100
}'
```
4) Create an account for the user 2
```shell
curl --header 'Authorization: Bearer token_user_2' --request POST 'http://localhost:8000/accounts'
```

2) Check the balance for the second account
```shell
curl --header 'Authorization: Bearer token_user_2' 'http://localhost:8000/accounts/2'
```

3) Transfer money between accounts
```shell
curl --request POST 'http://localhost:8000/transfer' \
--header 'Authorization: Bearer token_user_1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from": 1,
    "to": 2,
    "amount": 50
}'
```
