# gRPC User Microservice

A small, fully testable gRPC user service in Go: protobuf API, in-memory store, server/client binaries, and unit + in-process gRPC tests (`bufconn`).

## Layout

```
api/user/v1/user.proto     # service contract
gen/user/v1/               # generated protobuf / gRPC stubs
internal/store/            # in-memory persistence + unit tests
internal/service/          # gRPC handlers + bufconn tests
internal/gateway/          # HTTP JSON facade for browser testing
web/                       # API console UI
cmd/server/                # gRPC + web console entrypoint
cmd/client/                # CLI client for manual checks
scripts/gen.bat            # regenerate stubs
```

## Prerequisites

- Go 1.20+
- `protoc` (scripts expect `%LOCALAPPDATA%\protoc\bin\protoc.exe`)
- plugins: `protoc-gen-go`, `protoc-gen-go-grpc`

```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
```

## Generate stubs

```powershell
.\scripts\gen.bat
```

## Run tests

```powershell
go test ./...
```

## Run server / client

```powershell
go run ./cmd/server -addr :50051 -http :8080
```

Open the browser console: [http://127.0.0.1:8080](http://127.0.0.1:8080)

CLI client:

```powershell
go run ./cmd/client -action demo
```

Client actions: `demo`, `create`, `get`, `list`, `update`, `delete`.

### Browser HTTP mapping

| RPC | HTTP |
|-----|------|
| `CreateUser` | `POST /api/users` |
| `GetUser` | `GET /api/users/{id}` |
| `ListUsers` | `GET /api/users?page_size=&page_token=` |
| `UpdateUser` | `PUT /api/users/{id}` |
| `DeleteUser` | `DELETE /api/users/{id}` |

## API

| RPC | Behavior |
|-----|----------|
| `CreateUser` | Create by name/email; duplicate email → `AlreadyExists` |
| `GetUser` | Fetch by id; missing → `NotFound` |
| `ListUsers` | Cursor pagination via `page_token` |
| `UpdateUser` | Replace name/email |
| `DeleteUser` | Remove by id |

Server reflection is enabled for tools such as `grpcurl`.
