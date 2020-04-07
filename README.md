# GuruBackendService
This is backend service that can processes player transactions.

To create new user in a system use:
```
curl -X POST \
  http://localhost:8080/api/v1/user/create \
  -H 'Content-Type: application/json' \
  -d '{
	"id": 2,
	"balance": 0.0,
	"token": "testtest"
}'
```

To get user by ID use, where 'id' is a user id in a system.
```
curl -X POST \
  http://localhost:8080/api/v1/user/get \
  -H 'Content-Type: application/json' \
  -d '{
	"id": 2,
	"token": "testtest"
}'
```
To add deposit on user's account use, where 'id' is deposit id.
```
curl -X POST \
  http://localhost:8080/api/v1/user/deposit \
  -H 'Content-Type: application/json' \
  -d '{
	"id": 1,
	"token": "testtest",
	"userId": 2,
	"amount": 50
}'
```
To make transaction use curl below, where 'id' is transaction id, betType can be Win or Bet.
```
curl -X POST \
  http://localhost:8080/api/v1/transaction \
  -H 'Content-Type: application/json' \
  -d '{
	"id": 1,
	"token": "testtest",
	"userId": 2,
	"amount": 50,
	"betType": "Bet"
}'
```