# AP2 Assignment 2 — gRPC Microservices

**Student:** Syndaly Yerzhan, SE-2424

## Overview

This project implements a Medical Scheduling Platform using **gRPC** as the sole communication protocol, migrated from the REST-based Assignment 1.

The system consists of two independent services:

- **Doctor Service (port 50051)** — manages doctor data
- **Appointment Service (port 50052)** — manages appointments and depends on Doctor Service

---

### What changed from Assignment 1

| Feature | Assignment 1 | Assignment 2 |
|---|---|---|
| Communication | REST / HTTP | gRPC |
| Serialization format | JSON | Protocol Buffers (binary) |
| Transport framework | Gin | gRPC server |
| Inter-service calls | HTTP client | gRPC client stub |
| API contract | Implicit (no schema) | Strict `.proto` files |
| Error handling | HTTP status codes | gRPC status codes |

---

### What did NOT change

- domain models
- use-case logic
- business rules
- repository layer
- Clean Architecture layering

---

## Architecture

The project follows Clean Architecture with strict dependency direction:

```text
┌─────────────────────────────────────────────────────────┐
│                    Client / Postman                     │
└────────────────────────┬────────────────────────────────┘
                         │ gRPC
          ┌──────────────▼──────────────┐
          │    Transport Layer (gRPC)   │  ← thin handler, maps proto ↔ domain
          └──────────────┬──────────────┘
                         │ calls interface
          ┌──────────────▼──────────────┐
          │      Use Case Layer         │  ← all business rules live here
          └──────────────┬──────────────┘
                         │ calls interface
          ┌──────────────▼──────────────┐
          │     Repository Layer        │  ← in-memory storage
          └─────────────────────────────┘
```

---

## Appointment Service also holds a DoctorClient:

Appointment UseCase
        │
        │ DoctorClient interface
        ▼
gRPC Client Stub → Doctor Service (:50051)

---

## Key rules

- gRPC handlers are thin: unmarshal proto → call use case → marshal response
- Use cases never import protobuf types
- Proto ↔ domain mapping happens only in the transport layer
- DoctorClient is hidden behind an interface injected into the use case

---

## Project Structure

``` 
ap2-assignment2/
│
├── go.mod
├── README.md
│
├── doctor-service/
│   ├── cmd/doctor-service/main.go
│   ├── internal/
│   │   ├── app/app.go                        ← wires all dependencies
│   │   ├── model/doctor.go                   ← pure domain model
│   │   ├── repository/inmemory_doctor_repository.go
│   │   ├── transport/grpc/doctor_server.go   ← gRPC handler (thin)
│   │   └── usecase/doctor_usecase.go         ← business logic
│   └── proto/
│       ├── doctor.proto
│       ├── doctor.pb.go                      ← generated
│       └── doctor_grpc.pb.go                 ← generated
│
├── appointment-service/
│   ├── cmd/appointment-service/main.go
│   ├── internal/
│   │   ├── app/app.go
│   │   ├── client/doctor_client.go           ← gRPC client to Doctor Service
│   │   ├── model/appointment.go
│   │   ├── repository/inmemory_appointment_repository.go
│   │   ├── transport/grpc/appointment_server.go
│   │   └── usecase/appointment_usecase.go
│   └── proto/
│       ├── appointment.proto
│       ├── appointment.pb.go
│       └── appointment_grpc.pb.go
```

---

## How to Install protoc and Plugins

1. Install protoc compiler

- macOS: `brew install protobuf `
- Linux: `apt install -y protobuf-compiler`
- Verify: `protoc --version`

2. Install Go plugins

- `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

Make sure $GOPATH/bin is in your PATH:

- `export PATH="$PATH:$(go env GOPATH)/bin"`

---

## How to Regenerate Proto Stubs

Run from the project root (ap2-assignment2/):

---

## Doctor Service

`protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       doctor-service/proto/doctor.proto`

---

## Appointment Service

`protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       appointment-service/proto/appointment.proto`

---

### Generated files (*.pb.go, *_grpc.pb.go) are committed to the repository — no need to regenerate to run the project.

---

## How to Run

## Important: start Doctor Service first, then Appointment Service.
### Step 1 — Start Doctor Service

go run ./doctor-service/cmd/doctor-service
### Output:
Doctor gRPC service started at :50051

### Step 2 — Start Appointment Service
go run ./appointment-service/cmd/appointment-service
### Output:
Appointment gRPC service started at :50052

### Both services use in-memory storage — no database setup required.

--- 

## Service Responsibilities

### Doctor Service (:50051)
- Owns all doctor data
- Enforces email uniqueness
- Responds to doctor lookup requests from Appointment Service
### Appointment Service (:50052)
- Owns all appointment data
- Validates doctor existence by calling Doctor Service via gRPC before creating an appointment
- Manages appointment status lifecycle

--- 

## Proto Contract — RPC Description

### Doctor Service
| RPC          | Request                          | Response            | Business Rule                                      |
|--------------|----------------------------------|---------------------|----------------------------------------------------|
| CreateDoctor | full_name, specialization, email | DoctorResponse      | full_name and email required; email must be unique |
| GetDoctor    | id                               | DoctorResponse      | Returns NOT_FOUND if ID does not exist             |
| ListDoctors  | (empty)                          | ListDoctorsResponse | Returns all stored doctors                         |


### Appointment Service
| RPC                     | Request                       | Response                 | Business Rule                                                    |
|--------------------------|-------------------------------|--------------------------|------------------------------------------------------------------|
| CreateAppointment        | title, description, doctor_id | AppointmentResponse      | title and doctor_id required; doctor must exist (gRPC call)      |
| GetAppointment           | id                            | AppointmentResponse      | Returns NOT_FOUND if ID does not exist                           |
| ListAppointments         | (empty)                       | ListAppointmentsResponse | Returns all stored appointments                                  |
| UpdateAppointmentStatus  | id, status                    | AppointmentResponse      | Status must be new / in_progress / done; done → new is forbidden |

## Inter-Service Communication

When `CreateAppointment` or `UpdateAppointmentStatus` is called:

- Appointment Service use case calls `DoctorClient.GetDoctor(ctx, doctorID)`
- DoctorClient (in `internal/client/`) sends a gRPC request to Doctor Service at `localhost:50051`

Doctor Service responds:

- OK → doctor exists, proceed  
- NOT_FOUND → mapped to FailedPrecondition (doctor does not exist)  
- connection error → mapped to Unavailable (service is down)  

The DoctorClient is hidden behind an interface defined in the use case package — the use case does not know it is talking to gRPC.

## Error Handling

| Situation | gRPC Status Code |
|----------|------------------|
| Required field missing | INVALID_ARGUMENT |
| Email already in use | ALREADY_EXISTS |
| Doctor/Appointment ID not found (local) | NOT_FOUND |
| Doctor does not exist (remote check) | FAILED_PRECONDITION |
| Doctor Service unreachable | UNAVAILABLE |
| Invalid status transition (done → new) | INVALID_ARGUMENT |

---

## Failure Scenario

What happens when Doctor Service is unavailable:

- Appointment Service calls `DoctorClient.GetDoctor()`
- gRPC connection fails with a transport error
- DoctorClient catches the error and returns `ErrDoctorUnavailable`
- Use case returns the error to the handler
- Handler returns:

```go
status.Error(codes.Unavailable, "doctor service unavailable")
```

- Client receives a proper gRPC UNAVAILABLE status

In production, this would be extended with:

- timeouts on gRPC calls (context with deadline)
- retry with exponential backoff
- circuit breaker (e.g. using go-resilience or sony/gobreaker)
- health checks between services

## REST vs gRPC — Trade-off Discussion

| Aspect          | REST                           | gRPC                                               |
|-----------------|--------------------------------|----------------------------------------------------|
| Format          | JSON (text, human-readable)    | Protobuf (binary, compact)                         |
| Performance     | Slower — JSON parsing overhead | Faster — binary serialization, HTTP/2 multiplexing |
| Contract        | Implicit — no enforced schema  | Strict — .proto file defines exact types           |
| Code generation | Manual client code             | Auto-generated stubs from .proto                   |
| Browser support | Native                         | Requires gRPC-Web proxy                            |
| Best for        | Public APIs, browser clients   | Internal microservice communication                |

## Testing with grpcurl

Install grpcurl:

```bash
brew install grpcurl
```

## Doctor Service
Create a doctor:

```bash
grpcurl -plaintext -d '{
  "full_name": "Alice Smith",
  "specialization": "Cardiology",
  "email": "alice@hospital.com"
}' localhost:50051 doctor.DoctorService/CreateDoctor
```
Get doctor by ID:

```bash
grpcurl -plaintext -d '{"id": "<doctor-id>"}' \
  localhost:50051 doctor.DoctorService/GetDoctor
```
List all doctors:

```bash
grpcurl -plaintext -d '{}' localhost:50051 doctor.DoctorService/ListDoctors
```
## Appointment Service
Create an appointment:

```bash
grpcurl -plaintext -d '{
  "title": "Annual Checkup",
  "description": "Routine visit",
  "doctor_id": "<doctor-id>"
}' localhost:50052 appointment.AppointmentService/CreateAppointment
```
Get appointment by ID:

```bash
grpcurl -plaintext -d '{"id": "<appointment-id>"}' \
  localhost:50052 appointment.AppointmentService/GetAppointment
```
List all appointments:
```bash
grpcurl -plaintext -d '{}' localhost:50052 appointment.AppointmentService/ListAppointments
```
Update appointment status:

```bash
grpcurl -plaintext -d '{
  "id": "<appointment-id>",
  "status": "in_progress"
}' localhost:50052 appointment.AppointmentService/UpdateAppointmentStatus
```
Test error — invalid doctor ID:

```bash
grpcurl -plaintext -d '{
  "title": "Test",
  "doctor_id": "nonexistent-id"
}' localhost:50052 appointment.AppointmentService/CreateAppointment
```
Expected: `FAILED_PRECONDITION`

Test error — Doctor Service down:

```bash
grpcurl -plaintext -d '{
  "title": "Test",
  "doctor_id": "some-id"
}' localhost:50052 appointment.AppointmentService/CreateAppointment
```
Expected: `UNAVAILABLE`

## Improvements from Assignment 1

- REST → gRPC migration (transport layer fully replaced)
- Strict API contracts via .proto files
- Generated client/server stubs (no manual HTTP parsing)
- Domain models have no JSON tags (transport concern removed from domain)
- Handlers depend on interfaces, not concrete use case structs
- Use case owns its ports (repository and client interfaces defined in use case package)
- Proper gRPC status codes on all error paths


## Conclusion
This project demonstrates a successful migration from REST to gRPC while preserving Clean Architecture principles.
The system benefits from strict API contracts defined via Protocol Buffers, improved performance through binary serialization, and clear separation of concerns between layers.
Additionally, inter-service communication is now strongly typed and more efficient, making the architecture more scalable and production-ready.