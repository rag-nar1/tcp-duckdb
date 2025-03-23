# TCP Server for DuckDB

A TCP server that enables management and interaction with DuckDB databases over a network. The server provides functionality for user authentication, database operations, and access control.

## Table of Contents
- [Features](#features)
- [Commands Reference](#commands-reference)
  - [Login](#login)
  - [Create](#create)
  - [Connect](#connect)
  - [Grant](#grant)
  - [Query](#query)
  - [Link and Migrate](#link-and-migrate)
- [Usage Example](#usage-example)
- [Security and Access Control](#security-and-access-control)

## Features

- User authentication and authorization
- Database creation and management
- Table-level access control
- PostgreSQL database linking
- Transaction support
- Query execution

## Commands Reference

### Login

Authenticate to access the server.

```bash
login [username] [password]
```
- `[username]`: Your username
- `[password]`: Your password

**Example:**
```bash
login duck superpassword
```

> ⚠️ **IMPORTANT:** The super user is `duck`, which has privileges to create databases and users. The default password is `duck` - it is crucial to change this password immediately after setting up your project for security purposes.

### Create

Create databases or users (requires super user privileges).

#### Create Database
```bash
create database [database_name]
```

**Example:**
```bash
create database mydb
```

#### Create User
```bash
create user [username] [password]
```

**Example:**
```bash
create user john pass123
```

### Connect

Connect to a database to execute queries.

```bash
connect [database_name]
```

**Example:**
```bash
connect mydb
```

After connecting, you can:
- Execute single queries
- Start transactions with `start transaction`
- Commit with `commit`
- Rollback with `rollback`

### Grant

Grant access permissions to users (requires super user privileges).

#### Database Access
```bash
grant database [database_name] [username] [access_type]
```
Access types: `read`, `write`

**Example:**
```bash
grant database mydb john read
```

#### Table Access
```bash
grant table [database_name] [table_name] [username] [access_type...]
```
Access types: `select`, `update`, `insert`, `delete`

**Example:**
```bash
grant table mydb users john select insert
```

### Query

After connecting to a database, you can execute:

#### Single Query
```sql
SELECT * FROM table;
```

#### Transaction
```sql
start transaction
INSERT INTO users VALUES (1, 'John');
UPDATE users SET name = 'Johnny' WHERE id = 1;
commit
```

### Link and Migrate

Link DuckDB with PostgreSQL databases (requires super user privileges).

#### Link
```bash
link [database_name] [postgresql_connection_string]
```

**Example:**
```bash
link mydb "postgresql://user:password@localhost:5432/pgdb"
```

#### Migrate
```bash
migrate [database_name]
```

**Example:**
```bash
migrate mydb
```

**Note:** The `link` command establishes a connection between DuckDB and PostgreSQL by reading the PostgreSQL table schemas and recreating them in DuckDB, then copying all data from PostgreSQL tables into DuckDB. The `migrate` command maintains synchronization by reading the audit table to keep the DuckDB database in sync with PostgreSQL changes.

## Usage Example

Here's a step-by-step example of using the server:

1. **Login as Super User:**
```bash
login duck superpassword
```

2. **Create a Database:**
```bash
create database mydb
```

3. **Create a Regular User:**
```bash
create user analyst pass123
```

4. **Grant Database Access:**
```bash
grant database mydb analyst read
```

5. **Grant Table Access:**
```bash
grant table mydb customers analyst select
```

6. **User Login and Query:**
```bash
login analyst pass123
connect mydb
SELECT * FROM customers;
```

## Security and Access Control

The server enforces strict access control:
- Only the super user `duck` can create databases and users or grant permissions.
- Users can only connect to databases and execute queries on tables they have been granted access to.
- Permissions are checked for every connection and query operation.

This documentation provides the essentials for using the TCP Server for DuckDB. For further details, such as starting the server or handling advanced features, refer to the project's additional documentation or source code.