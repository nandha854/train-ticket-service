# Train Ticket Service

This service manages train ticket bookings, including seat assignments, ticket purchases, and modifications.

## Features

- **Seat Management**: Assign, release, and modify seats.
- **Ticket Management**: Purchase tickets, retrieve receipts, and manage user tickets.

## gRPC API Specifications

### PurchaseTicket

**RPC**: `PurchaseTicket`

**Request**:
```protobuf
message PurchaseTicketRequest {
  User user = 1;
  string from = 2;
  string to = 3;
}
```

**Response**:
```protobuf
message TicketReceipt {
  User user = 1;
  string from = 2;
  string to = 3;
  int32 price = 4;
  Seat seat = 5;
}
```

### GetReceipt

**RPC**: `GetReceipt`

**Request**:
```protobuf
message GetReceiptRequest {
  string email = 1;
}
```

**Response**:
```protobuf
message TicketReceipt {
  User user = 1;
  string from = 2;
  string to = 3;
  int32 price = 4;
  Seat seat = 5;
}
```

### GetUsersBySection

**RPC**: `GetUsersBySection`

**Request**:
```protobuf
message GetUsersBySectionRequest {
  string section = 1;
}
```seat

**Response**:
```protobuf
message UsersBySectionResponse {
  repeated UserTicket users = 1;
}
```

### RemoveUser

**RPC**: `RemoveUser`

**Request**:
```protobuf
message RemoveUserRequest {
  string email = 1;
}
```

**Response**:
```protobuf
message RemoveUserResponse {
  string message = 1;
}
```

### ModifyUserSeat

**RPC**: `ModifyUserSeat`

**Request**:
```protobuf
message ModifyUserSeatRequest {
  string email = 1;
  Seat newSeat = 2;
}
```

**Response**:
```protobuf
message TicketReceipt {
  User user = 1;
  string from = 2;
  string to = 3;
  int32 price = 4;
  Seat seat = 5;
}
```

## Running Tests

To run the tests for the SeatManager, use the following command:

```sh
go test ./service -v
```