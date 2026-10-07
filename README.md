# Pismo Tech Case — Transactions routine

A small HTTP API in Go that manages cardholder accounts and their transactions.

## Running

The easiest way is the `./run` script (requires Docker):

```sh
./run         # build and start the API on http://localhost:8080 (same as ./run up)
./run test    # run the test suite (race detector + coverage) in a Go container
./run down    # stop and remove the containers
```

The port can be changed with the `PORT` environment variable (default `8080`):

```sh
PORT=9090 ./run up
```

### Without the script

With Docker Compose:

```sh
docker compose up --build
```

Locally, requires Go 1.27.1+:

```sh
go run .            # listens on :8080
PORT=9090 go run .  # custom port
go test ./...       # tests, without Docker
```

## API

| Method | Path                     | Description               |
|--------|--------------------------|---------------------------|
| POST   | `/accounts`              | Create an account         |
| GET    | `/accounts/{account_id}` | Retrieve an account       |
| POST   | `/transactions`          | Create a transaction      |
| GET    | `/health`                | Health check              |

All errors are returned as JSON: `{"error": "<message>"}`.

### Create account

```sh
curl -X POST http://localhost:8080/accounts \
  -H 'Content-Type: application/json' \
  -d '{"document_number": "12345678900"}'
```

```json
{"account_id": 1, "document_number": "12345678900"}
```

| Status | When                          |
|--------|-------------------------------|
| 201    | Account created               |
| 400    | Invalid JSON / missing `document_number` |

### Get account

```sh
curl http://localhost:8080/accounts/1
```

```json
{"account_id": 1, "document_number": "12345678900"}
```

| Status | When                                  |
|--------|---------------------------------------|
| 200    | Account found                         |
| 400    | `account_id` is not a positive integer |
| 404    | Account does not exist                |

### Create transaction

```sh
curl -X POST http://localhost:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{"account_id": 1, "operation_type_id": 4, "amount": 123.45}'
```

```json
{
  "transaction_id": 1,
  "account_id": 1,
  "operation_type_id": 4,
  "amount": 123.45,
  "event_date": "2026-10-07T10:49:34.505747841Z"
}
```

| Status | When                                                        |
|--------|-------------------------------------------------------------|
| 201    | Transaction created                                         |
| 400    | Invalid JSON, missing/zero field, invalid operation type, more than 2 decimal places |
| 422    | `account_id` does not reference an existing account         |

#### Operation types

| ID | Description                | Stored sign |
|----|----------------------------|-------------|
| 1  | Normal Purchase            | negative    |
| 2  | Purchase with installments | negative    |
| 3  | Withdrawal                 | negative    |
| 4  | Credit Voucher             | positive    |

## Project structure

```
.
├── main.go          # entrypoint + wiring (repositories → services → API)
├── api.go           # HTTP handlers, request decoding, error → status mapping
├── account/         # Account model, service, repository interface + in-memory impl
├── transaction/     # Transaction model, operation types, service, repository
├── run              # helper script: up / test / down
├── Dockerfile       # multi-stage build, distroless runtime image
└── compose.yml      # runs the API (port set by PORT, default 8080)
```

Each domain package owns its model, its business rules (service) and a `Repository` interface. The HTTP layer only translates between HTTP and the services; it holds no business rules.

## Design decisions

### Amount sign is defined by the operation type

The client may send the amount with any sign; the server normalizes it:
debit operations (purchases and withdrawals) are stored as negative and credit vouchers as positive. Zero is rejected.

This rule lives in `transaction.NewTransaction`, so it is impossible to create a transaction with an inconsistent sign anywhere in the code.

*Trade-off:* rejecting negative input would be stricter and would surface client bugs earlier. Normalizing was chosen because the operation type already carries the intent, so the sign from the client is redundant information.

### Money as `decimal`, not `float64`

Amounts use `shopspring/decimal` to avoid floating-point rounding errors.
Amounts with more than 2 decimal places are rejected instead of silently rounded.

### Operation types as Go constants

Operation types are a typed enum (`transaction.OperationType`) with behavior (`IsValid`, `IsDebit`). They change rarely and drive business logic, so keeping them in code makes the rules explicit and testable. With a real database they could also be a table, used only for referential integrity.

### Non-existent account returns 422

The request is syntactically valid, but refers to an account that does not exist, so `422 Unprocessable Entity` is used instead of `404` (the `/transactions` resource itself was found).

### `event_date` is set by the server

The event date is generated in UTC by the service. The clock is injected (`func() time.Time`) so tests are deterministic.

### In-memory persistence

Repositories are in-memory (map + mutex) behind an interface. Swapping to a real database (e.g. Postgres) means adding a new `Repository` implementation and changing only the wiring in `main.go`.

## Testing

- **Unit tests** for models, services and repositories.
- **Handler tests** for the HTTP layer (status codes and error mapping).
- **Integration test** (`api_integration_test.go`) exercising the full wiring through HTTP.

Two styles of test doubles are used on purpose: `testify/mock` mocks (to assert interactions and simulate failures) and the real in-memory repository (to test behavior end to end). Similarly, some tests are table-driven and others are standalone, depending on what reads best for each case.

## Possible next steps

- Real database (Postgres) with migrations.
- `GET /transactions` / account balance.
- Structured logging and request IDs.
