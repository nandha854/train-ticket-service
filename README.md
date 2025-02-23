# Train Ticket Booking Service - gRPC API

## Overview
The **Train Ticket Booking Service** is a gRPC-based system that allows users to book tickets, retrieve receipts, manage seat assignments, and cancel bookings.

## gRPC Service Definition
The service is defined in the **ticketBooking.proto** file and includes the following RPC methods:

### **TicketService**
```proto
service TicketService {
  rpc PurchaseTicket(PurchaseTicketRequest) returns (TicketReceipt) {}
  rpc GetReceipt(GetReceiptRequest) returns (TicketReceipt) {}
  rpc GetUsersBySection(GetUsersBySectionRequest) returns (UsersBySectionResponse) {}
  rpc RemoveUser(RemoveUserRequest) returns (RemoveUserResponse) {}
  rpc ModifyUserSeat(ModifyUserSeatRequest) returns (TicketReceipt) {}
}
```

## Features
### **1. Ticket Management**
- **PurchaseTicket:** Allows users to purchase tickets and assigns them a seat.
- **GetReceipt:** Retrieves the ticket receipt for a specific user.
- **RemoveUser:** Cancels a ticket and releases the assigned seat.
- **ModifyUserSeat:** Allows users to change their seat allocation.

### **2. Seat Management**
- **Seat allocation:** Seats are assigned in a round-robin manner across sections.
- **Seat modification:** Users can request to change their assigned seats.
- **Seat release:** When a ticket is canceled, the seat becomes available again.
- **Section-based queries:** Retrieve users seated in a specific section.

## Messages Definition

### **User Information**
```proto
message User {
  string first_name = 1;
  string last_name = 2;
  string email = 3;
}
```

### **Ticket Booking Requests & Responses**
```proto
message PurchaseTicketRequest {
  string from = 1;
  string to = 2;
  User user = 3;
}

message TicketReceipt {
  string from = 1;
  string to = 2;
  User user = 3;
  double price = 4;
  Seat seat = 5;
}
```

### **Seat Management**
```proto
message Seat {
  string section = 1;
  int32 seat_number = 2;
}
```

### **Ticket Lookup & Cancellation**
```proto
message GetReceiptRequest {
  string email = 1;
}

message RemoveUserRequest {
  string email = 1;
}

message RemoveUserResponse {
  string message = 1;
}
```

### **Section-wise User Retrieval**
```proto
message GetUsersBySectionRequest {
  string section = 1;
}

message UserTicket {
  User user = 3;
  Seat seat = 5;
}

message UsersBySectionResponse {
  repeated UserTicket users = 1;
}
```

### **Seat Modification**
```proto
message ModifyUserSeatRequest {
  string email = 1;
  Seat new_seat = 2;
}
```

## Running the Service
### **1. Install Dependencies**
Ensure you have `protoc` installed and the Go plugins for gRPC:
```sh
brew install protobuf  # macOS
sudo apt install protobuf-compiler  # Linux
```

Install gRPC dependencies for Go:
```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### **2. Compile the Protocol Buffers**
```sh
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/ticketBooking.proto
```

### **3. Run the Server**
```sh
go run main.go
```

### **4. Client Request Example**
You can use a clients/examples.go to play with server
```sh
go run client/examples.go
```

