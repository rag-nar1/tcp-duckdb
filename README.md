# TCP-DuckDB Documentation

## Table of Contents

- [Introduction](#introduction)
  - [Key Features](#key-features)
- [Setup Guide](#setup-guide)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Tools and Technologies](#tools-and-technologies)
  - [Programming Languages](#programming-languages)
  - [Database Technologies](#database-technologies)
  - [Libraries and Frameworks](#libraries-and-frameworks)
  - [Development Tools](#development-tools)
- [Command Reference](#command-reference)
  - [1. Login Command](#1-login-command)
  - [2. Create Command](#2-create-command)
    - [Create Database](#create-database)
    - [Create User](#create-user)
  - [3. Connect Command](#3-connect-command)
  - [4. Grant Command](#4-grant-command)
    - [Grant Database Access](#grant-database-access)
    - [Grant Table Access](#grant-table-access)
  - [5. Link Command](#5-link-command)
  - [6. Migrate Command](#6-migrate-command)
  - [7. Update Command](#7-update-command)
    - [Update Database Name](#update-database-name)
    - [Update Username](#update-username)
    - [Update User Password](#update-user-password)
- [Transaction Management](#transaction-management)
- [Error Handling](#error-handling)
- [Environment Configuration](#environment-configuration)
- [Deployment](#deployment)
- [Internal Architecture](#internal-architecture)
  - [Core Components and Modules](#core-components-and-modules)
    - [1. Server Module](#1-server-module-server)
    - [2. Main Module](#2-main-module-main)
    - [3. Request Handler Module](#3-request-handler-module-request_handler)
    - [4. Connection Pool Module](#4-connection-pool-module-pool)
    - [5. Login Module](#5-login-module-login)
    - [6. Create Module](#6-create-module-create)
    - [7. Connect Module](#7-connect-module-connect)
    - [8. Grant Module](#8-grant-module-grant)
    - [9. Link Module](#9-link-module-link)
    - [10. Migrate Module](#10-migrate-module-migrate)
    - [11. Update Module](#11-update-module-update)
    - [12. Utils Module](#12-utils-module-utils)
  - [Request Processing Flow](#request-processing-flow)
- [Todo](#todo)
- [License](#license)

## Introduction

TCP-DuckDB is a TCP server implementation that provides networked access to DuckDB databases. The server enables remote database management with features like user authentication, access control, and PostgreSQL integration. Written in Go, it leverages the power of DuckDB, a lightweight analytical database engine.

### Key Features

- **TCP Interface**: Network-accessible database service
- **User Authentication**: Multi-user support with authentication
- **Database Management**: Create and manage DuckDB databases
- **Permission Control**: Fine-grained access permissions at database and table levels
- **PostgreSQL Integration**: Link with PostgreSQL databases
- **Transaction Support**: Full transaction support for data operations
- **Connection Pooling**: Efficient database connection management

## Setup Guide

### Installation

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed on your system
- Git to clone this repository

### Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/rag-nar1/tcp-duckdb.git
   cd TCP-Duckdb
   ```

2. Build and start the server using Docker Compose:
   ```bash
   docker compose up --build
   ```

3. Check logs to verify the server is running:
   ```bash
   docker logs tcp-duckdb-tcp-duckdb-1
   ```
   
   You should see output like:
   ```
   INFO    YYYY/MM/DD HH:MM:SS Super user created
   INFO    YYYY/MM/DD HH:MM:SS listening to 0.0.0.0:4000
   ```

4. Stop the server when finished:
   ```bash
   docker compose down
   ```

### Manual Build

If you prefer to build and run manually:

1. Build the Docker image:
   ```bash
   docker build -t tcp-duckdb .
   ```

2. Run the container:
   ```bash
   docker run -d -p 4000:4000 \
     -v $(pwd)/storge:/app/storge \
     -v $(pwd)/server:/app/server \
     -e ServerPort=4000 \
     -e ServerAddr=0.0.0.0 \
     -e DBdir=/app/storge/server/ \
     -e ServerDbFile=db.sqlite3 \
     -e ENCRYPTION_KEY=A15pG0m3hwf0tfpVW6m92eZ6vRmAQA3C \
     --name tcp-duckdb-container \
     tcp-duckdb
   ```

### Configuration

The server can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| ServerPort | Port the server listens on | 4000 |
| ServerAddr | Address the server binds to | 0.0.0.0 |
| DBdir | Directory for the SQLite database | /app/storge/server/ |
| ServerDbFile | SQLite database filename | db.sqlite3 |
| ENCRYPTION_KEY | Key used for encryption | `ENCRYPTION_KEY` |

You can modify these values in the `docker-compose.yml` file or pass them directly when running the container.

### Development

To build and run the application locally:

1. Install Go 1.24 or later
2. Install SQLite development libraries
3. Clone the repository
4. Run:
   ```bash
   go mod download
   go build -o ./build/server main/*
   ./build/server
   ```

### Troubleshooting

**Database initialization errors**

If you see errors related to database tables, ensure the schema is correctly applied:

```bash
# Connect to the running container
docker exec -it tcp-duckdb-tcp-duckdb-1 bash

# Verify the database exists
ls -la /app/storge/server/

# Manually apply the schema if needed
sqlite3 /app/storge/server/db.sqlite3 < /app/storge/server/scheme.sql
```

**Connection issues**

The server listens on TCP port 4000. Verify the port is correctly mapped and not blocked by a firewall.

1. Clone the repository:
```bash
git clone https://github.com/rag-nar1/tcp-duckdb.git
cd TCP-Duckdb
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
Create or modify the `.env` file in the project root:
```env
ServerPort=4000
ServerAddr=localhost
DBdir=/path/to/storage/
ServerDbFile=server/db.sqlite3
ENCRYPTION_KEY="YourEncryptionKey"
```

4. Build the project:
```bash
make build
```

5. Run the server:
```bash
make run
```

## Tools and Technologies

### Programming Languages

- **Go (Golang)**: The primary programming language used for the entire codebase. Go was chosen for its efficiency in building networked services, excellent concurrency support through goroutines, and strong standard library.

### Database Technologies

- **DuckDB**: A lightweight, in-process analytical database management system. It serves as the primary storage engine for the application, providing fast analytical query capabilities.

- **SQLite**: Used for storing user authentication and permission data. SQLite was chosen for its simplicity, reliability, and zero-configuration nature.

- **PostgreSQL**: Supported as an optional integration, allowing linking and synchronization with PostgreSQL databases. The system can replicate schema and data from PostgreSQL into DuckDB.

### Libraries and Frameworks

- **go-duckdb**: The Go driver for DuckDB that enables interaction with DuckDB databases from Go code.

- **go-sqlite3**: The Go interface to the SQLite3 database, used for user management.

- **lib/pq**: PostgreSQL driver for Go, used for connecting to PostgreSQL databases when using the link functionality.

- **godotenv**: Used for loading environment variables from .env files.

- **Standard Library Packages**: 
  - `net`: Core networking functionality for TCP server implementation
  - `database/sql`: Database interaction
  - `sync`: Synchronization primitives for concurrent operations
  - `bufio`: Buffered I/O operations
  - `crypto`: Cryptographic functions for secure password hashing

### Development Tools

- **Makefile**: Used for build automation, with predefined tasks for building, running, and code formatting.

- **Git**: Version control system with custom pre-commit hooks for code formatting.

- **Environment Configuration**: Uses .env files for configuration management.

- **Connection Pool Implementation**: Custom LRU (Least Recently Used) cache implementation for efficient database connection management.

## Command Reference

### 1. Login Command

**[All Users]**

Authenticates a user to access the server. This is the first command that must be executed before any other operation can be performed.

```
login [username] [password]
```

**Authentication Process:**
1. The client sends the login command with username and password
2. The server validates the credentials against the SQLite user database
3. If successful, a user session is established with appropriate privileges
4. All subsequent commands will operate under this authenticated user context

**Super User Information:**
- The default super user is `duck` with initial password `duck`
- The super user has full administrative privileges including:
  - Creating and managing databases
  - Creating and managing users
  - Granting permissions
  - Linking with PostgreSQL databases
  - Performing update operations
- For security reasons, it is strongly recommended to change the super user password after initial setup using the update command

**Example:**
```
login duck duck
```
Response on success:
```
success
```

### 2. Create Command

**[Super User Only]**

#### Create Database
Creates a new DuckDB database (requires super user privileges).

```
create database [database_name]
```

**Example:**
```
create database analytics
```
Response on success:
```
success
```

#### Create User
Creates a new user (requires super user privileges).

```
create user [username] [password]
```

**Example:**
```
create user analyst securepassword
```
Response on success:
```
success
```

### 3. Connect Command

**[All Users]**

Connects to an existing database to execute queries.

```
connect [database_name]
```

After connecting, the system:
- Verifies database existence
- Checks user permissions
- Acquires database connection from pool
- Allows executing queries or starting transactions

**Example:**
```
connect analytics
```
Response on success:
```
success
```

Once connected, you can execute SQL queries directly:
```
SELECT * FROM users;
```

### 4. Grant Command

**[Super User Only]**

#### Grant Database Access
Grants database access to a user (requires super user privileges).

```
grant database [database_name] [username] [access_type]
```

Access types:
- `read`: Read-only access
- `write`: Read and write access

**Example:**
```
grant database analytics analyst read
```
Response on success:
```
success
```

#### Grant Table Access
Grants table-level permissions to a user (requires super user privileges).

```
grant table [database_name] [table_name] [username] [access_type...]
```

Access types:
- `select`: Permission to query the table
- `insert`: Permission to insert data
- `update`: Permission to update data
- `delete`: Permission to delete data

**Example:**
```
grant table analytics users analyst select insert
```
Response on success:
```
success
```

### 5. Link Command

**[Super User Only]**

Links a DuckDB database with a PostgreSQL database (requires super user privileges).

```
link [database_name] [postgresql_connection_string]
```

Implementation:
- Connects to the PostgreSQL database
- Retrieves schema information
- Creates corresponding tables in DuckDB
- Copies data from PostgreSQL to DuckDB
- Sets up audit triggers for change tracking

**Example:**
```
link analytics "postgresql://user:password@localhost:5432/analytics_pg"
```
Response on success:
```
success
```

### 6. Migrate Command

**[Super User Only]**

Synchronizes changes from a linked PostgreSQL database to DuckDB (requires super user privileges).

```
migrate [database_name]
```

Implementation:
- Reads audit logs from PostgreSQL
- Applies changes to the DuckDB database
- Updates tracking information

**Example:**
```
migrate analytics
```
Response on success:
```
success
```

### 7. Update Command

**[Super User Only]**

Updates database names or user information (requires super user privileges).

The update command has three variations:

#### Update Database Name
```
update database [old_database_name] [new_database_name]
```

Implementation:
- Verifies database existence
- Renames the database file
- Updates database name in server records

**Example:**
```
update database analytics analytics_prod
```
Response on success:
```
success
```

#### Update Username
```
update user username [old_username] [new_username]
```

Implementation:
- Verifies user existence
- Updates username in user database

**Example:**
```
update user username analyst data_scientist
```
Response on success:
```
success
```

#### Update User Password
```
update user password [username] [new_password]
```

Implementation:
- Verifies user existence
- Hashes the new password
- Updates password in user database

**Example:**
```
update user password analyst new_secure_password
```
Response on success:
```
success
```

## Transaction Management

**[All Users with Database Access]**

After connecting to a database, you can manage transactions:

#### Start Transaction
```
start transaction
```

#### Execute Queries in Transaction
```
INSERT INTO users VALUES (1, 'John');
UPDATE users SET name = 'Johnny' WHERE id = 1;
```

#### Commit Transaction
```
commit
```

#### Rollback Transaction
```
rollback
```

**Example:**
```
connect analytics
start transaction
INSERT INTO users VALUES (1, 'John');
UPDATE users SET name = 'Johnny' WHERE id = 1;
commit
```

## Error Handling

The server implements structured error responses:
- `response.BadRequest(writer)`
- `response.InternalError(writer)`
- `response.UnauthorizedError(writer)`
- `response.DoesNotExistDatabse(writer, dbname)`
- `response.AccesDeniedOverDatabase(writer, UserName, dbname)`

## Environment Configuration

The server uses the following environment variables:
- `ServerPort`: TCP port for the server (default: 4000)
- `ServerAddr`: Server address (default: localhost)
- `DBdir`: Directory for storing databases
- `ServerDbFile`: Path to the server's SQLite database
- `ENCRYPTION_KEY`: Key for encrypting/decrypting PostgreSQL connection strings

## Deployment

The server can be deployed on any system with Go and DuckDB installed:

1. Build the server:
```bash
make build
```

2. Configure environment in `.env`

3. Run the server:
```bash
make run
```

## Internal Architecture

### Core Components and Modules

#### 1. Server Module (`server/`)

The main server component is responsible for the core server functionality:

- **Configuration Management**: Loads environment variables and configures the server
- **Database Connection**: Maintains connection to the SQLite user database
- **Statement Preparation**: Prepares SQL statements for efficient execution
- **Super User Management**: Creates and manages the super user account
- **Logging**: Implements structured logging for errors and information

**Key Files**:
- `config.go`: Server configuration and initialization
- `server.go`: Core server functionality implementation

**Key Functions**:
- `NewServer()`: Initializes server with configurations
- `CreateSuper()`: Creates the initial super user if it doesn't exist
- `PrepareStmt()`: Prepares SQL statements for later use

The Server struct centralizes server state:
```go
type Server struct {
    Sqlitedb *sql.DB          // SQLite database connection
    Dbstmt   map[string]*sql.Stmt  // Prepared statements
    Pool     *request_handler.RequestHandler  // Connection pool
    Port     string           // Server port
    Address  string           // Full server address
    InfoLog  *log.Logger      // Information logger
    ErrorLog *log.Logger      // Error logger
}
```

#### 2. Main Module (`main/`)

The entry point for the TCP server:

- **Server Initialization**: Initializes the server components
- **TCP Listener**: Sets up TCP socket and listens for connections
- **Connection Handling**: Accepts connections and spawns goroutines
- **Command Routing**: Routes incoming commands to appropriate handlers

**Key Files**:
- `main.go`: Entry point with TCP listener
- `router.go`: Command routing implementation

**Key Functions**:
- `main()`: Starts the TCP server listening on configured port
- `HandleConnection()`: Processes each client connection
- `Router()`: Routes requests to appropriate command handlers

#### 3. Request Handler Module (`request_handler/`)

Manages the lifecycle of database requests:

- **Request Queue**: Maintains queue of database connection requests
- **Connection Pooling Integration**: Works with pool module
- **Concurrency Management**: Handles simultaneous requests safely

**Key Files**:
- `request_handler.go`: Core request handling logic

**Key Functions**:
- `HandleRequest()`: Processes each database request
- `Spin()`: Starts the request handling background process
- `Push()`: Adds new requests to the queue

#### 4. Connection Pool Module (`pool/`)

Implements efficient connection management for DuckDB databases:

- **LRU Cache**: Implements Least Recently Used replacement policy
- **Connection Limits**: Manages maximum number of open connections
- **Pin Counting**: Tracks active database usage
- **Resource Management**: Efficiently manages database handles

**Key Files**:
- `pool.go`: Connection pool implementation
- `lru.go`: LRU cache implementation for connection eviction

**Key Functions**:
- `Get()`: Retrieves a database connection from the pool
- `NewPool()`: Creates a new connection pool
- `RecordAccess()`: Updates access time for LRU tracking

#### 5. Login Module (`login/`)

Handles user authentication:

- **Credential Verification**: Validates username and password
- **Password Hashing**: Securely stores and validates passwords
- **Session Establishment**: Sets up user session after authentication

**Key Files**:
- `handler.go`: Authentication request handling
- `service.go`: Authentication logic implementation

**Key Functions**:
- `Handler()`: Processes login requests
- `Login()`: Validates credentials against database

#### 6. Create Module (`create/`)

Manages creation of databases and users:

- **Database Creation**: Creates new DuckDB databases
- **User Creation**: Creates new user accounts
- **Permission Initialization**: Sets up initial permissions

**Key Files**:
- `handler.go`: Request handling for creation operations
- `service.go`: Implementation of creation operations

**Key Functions**:
- `CreateDatabase()`: Creates a new database
- `CreateUser()`: Creates a new user

#### 7. Connect Module (`connect/`)

Manages database connections and query execution:

- **Connection Establishment**: Connects to specified database
- **Permission Checking**: Verifies user has access to database
- **Query Execution**: Executes SQL queries on connected database
- **Transaction Management**: Handles SQL transactions

**Key Files**:
- `handler.go`: Connection request handling
- `service.go`: Query execution implementation
- `transaction.go`: Transaction management

**Key Functions**:
- `Handler()`: Processes connection requests
- `QueryService()`: Executes individual queries
- `Transaction()`: Manages database transactions

#### 8. Grant Module (`grant/`)

Manages access permissions:

- **Database Permissions**: Controls database access rights
- **Table Permissions**: Controls table-level access rights
- **Permission Checking**: Verifies permissions before granting

**Key Files**:
- `handler.go`: Grant request handling
- `service.go`: Permission management implementation

**Key Functions**:
- `GrantDatabaseAccess()`: Grants database-level access
- `GrantTableAccess()`: Grants table-level access

#### 9. Link Module (`link/`)

Facilitates PostgreSQL database integration:

- **Connection Management**: Establishes connections to PostgreSQL
- **Schema Transfer**: Replicates PostgreSQL schema to DuckDB
- **Data Migration**: Copies data from PostgreSQL to DuckDB
- **Connection String Encryption**: Securely stores PostgreSQL credentials

**Key Files**:
- `handler.go`: Link request handling
- `service.go`: Link implementation

**Key Functions**:
- `Link()`: Establishes connection and migrates schema/data

#### 10. Migrate Module (`migrate/`)

Handles data synchronization between PostgreSQL and DuckDB:

- **Change Detection**: Identifies changes in PostgreSQL
- **Synchronization**: Applies changes to DuckDB
- **Audit Log**: Processes audit logs for changes

**Key Files**:
- `handler.go`: Migrate request handling
- `service.go`: Migration implementation

**Key Functions**:
- `Migrate()`: Synchronizes changes from PostgreSQL to DuckDB

#### 11. Update Module (`update/`)

Manages updates to database and user information:

- **Database Name Updates**: Renames databases
- **User Information Updates**: Updates usernames and passwords
- **Validation**: Verifies existence before updates

**Key Files**:
- `handler.go`: Update request handling
- `service.go`: Update implementation

**Key Functions**:
- `UpdateDatabase()`: Renames a database
- `UpdateUserUsername()`: Updates a user's username
- `UpdateUserPassword()`: Updates a user's password

#### 12. Utils Module (`utils/`)

Provides utility functions used throughout the application:

- **Password Hashing**: Secures user passwords
- **Encryption**: Handles AES encryption for sensitive data
- **Path Management**: Manages file paths for databases
- **String Handling**: Provides string manipulation utilities

**Key Files**:
- `utils.go`: General utility functions
- `crypto.go`: Cryptographic functions

**Key Functions**:
- `Hash()`: Hashes passwords
- `Encrypt()/Decrypt()`: Encrypts/decrypts data with AES
- `UserDbPath()`: Resolves database file paths

### Request Processing Flow

1. **Connection Establishment**
   - Client connects to TCP server via `main.go`
   - Server spawns a goroutine for the connection via `HandleConnection()`
   - Client must authenticate via `login.Handler()`

2. **Command Processing**
   - After authentication, `Router()` in `main/router.go` processes requests
   - Commands are parsed and validated
   - Requests are routed to appropriate module handlers
   - Responses are sent back to client

3. **Database Operations**
   - Database connections are managed by the pool module
   - Operations are checked against user permissions
   - Transactions are handled with ACID guarantees
   - Results are returned to client

4. **PostgreSQL Integration Process**
   - Link operations copy schema and data from PostgreSQL
   - Migrate operations synchronize changes from PostgreSQL
   - Audit tables track changes for synchronization 

## Todo

This section outlines planned enhancements and improvements for the TCP-DuckDB project:

| Task | Description | Priority |
|------|-------------|----------|
| Client Library | Develop client libraries in multiple languages (Python, JavaScript, Java), go is in progress | High |
| Connection Encryption | Implement TLS/SSL for secure client-server communication | High |
| Change Data Capture | Swap the Audit table with CDC using [Debezium](https://debezium.io/) and [kafka](https://kafka.apache.org/) | Medium |
| Backup & Restore | Add automated backup and point-in-time recovery functionality | Medium |
| Query Caching | Add intelligent query result caching | Low |
| Web Admin Interface | Create a web-based administration interface | Low |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 