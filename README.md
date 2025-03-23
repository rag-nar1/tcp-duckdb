<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TCP Server for DuckDB Documentation</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            line-height: 1.6;
        }
        h1, h2, h3 {
            color: #333;
        }
        pre {
            background-color: #f4f4f4;
            padding: 10px;
            border-radius: 5px;
            overflow-x: auto;
        }
        code {
            font-family: monospace;
        }
        ul, ol {
            margin: 10px 0;
            padding-left: 20px;
        }
        a {
            color: #0066cc;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <h1>TCP Server for DuckDB</h1>
    <p>This TCP server enables management and interaction with DuckDB databases over a network. It provides a set of commands for user authentication, database and user creation, database connection, access control, and query execution. Below is a detailed guide on how to use the server with its supported commands.</p>

    <h2>Table of Contents</h2>
    <ul>
        <li><a href="#communication-protocol">Communication Protocol</a></li>
        <li><a href="#login">Login</a></li>
        <li><a href="#create">Create</a></li>
        <li><a href="#connect">Connect</a></li>
        <li><a href="#grant">Grant</a></li>
        <li><a href="#query">Query</a></li>
        <li><a href="#link-and-migrate">Link and Migrate</a></li>
        <li><a href="#usage-example">Usage Example</a></li>
        <li><a href="#security-and-access-control">Security and Access Control</a></li>
    </ul>

    <h2 id="communication-protocol">Communication Protocol</h2>
    <p>Commands are sent as strings over a TCP connection to the server, which processes them and returns appropriate responses. Each session typically begins with a login, followed by other commands based on user permissions.</p>

    <h2 id="login">Login</h2>
    <p>To interact with the server, you must first authenticate by logging in with a username and password.</p>
    <p><strong>Command:</strong></p>
    <pre><code>login [Username] [Password]</code></pre>
    <p>- <code>[Username]</code>: Your username.</p>
    <p>- <code>[Password]</code>: Your password.</p>
    <p><strong>Note:</strong> The super user is <code>duck</code>, which has privileges to create databases and users.</p>
    <p><strong>Example:</strong></p>
    <pre><code>login duck superpassword</code></pre>

    <h2 id="create">Create</h2>
    <p>The <code>create</code> command allows the super user to create new databases or users. You must be logged in as <code>duck</code> to use this command.</p>

    <h3>Create Database</h3>
    <p><strong>Command:</strong></p>
    <pre><code>create database [Database_name]</code></pre>
    <p>- <code>[Database_name]</code>: The name of the database to create. Must be unique across the server.</p>
    <p><strong>Example:</strong></p>
    <pre><code>create database mydb</code></pre>

    <h3>Create User</h3>
    <p><strong>Command:</strong></p>
    <pre><code>create user [Username] [Password]</code></pre>
    <p>- <code>[Username]</code>: The username for the new user. Must be unique across the server.</p>
    <p>- <code>[Password]</code>: The password for the new user.</p>
    <p><strong>Note:</strong> New users do not have access to any databases by default. Use the <code>grant</code> command to assign permissions.</p>
    <p><strong>Example:</strong></p>
    <pre><code>create user myuser 12345678</code></pre>

    <h2 id="connect">Connect</h2>
    <p>The <code>connect</code> command establishes a connection to a specific database. You must be logged in and have appropriate access permissions.</p>
    <p><strong>Command:</strong></p>
    <pre><code>connect [Database_name]</code></pre>
    <p>- <code>[Database_name]</code>: The name of the database to connect to.</p>
    <p><strong>Response:</strong> If successful, the server returns <code>"success"</code>.</p>
    <p><strong>Notes:</strong></p>
    <ul>
        <li>Requires prior login.</li>
        <li>You can only connect to databases you are authorized to access.</li>
    </ul>
    <p><strong>Example:</strong></p>
    <pre><code>connect mydb</code></pre>

    <h2 id="grant">Grant</h2>
    <p>The <code>grant</code> command, available to the super user, assigns access permissions to users for databases or tables.</p>

    <h3>Database-Level Access</h3>
    <p><strong>Command:</strong></p>
    <pre><code>grant database [Database_name] [username] [access type]</code></pre>
    <p>- <code>[Database_name]</code>: The database to grant access to.</p>
    <p>- <code>[username]</code>: The user receiving the access.</p>
    <p>- <code>[access type]</code>: Either <code>[read]</code> or <code>[write]</code>.</p>
    <p><strong>Example:</strong></p>
    <pre><code>grant database mydb myuser read</code></pre>

    <h3>Table-Level Access</h3>
    <p><strong>Command:</strong></p>
    <pre><code>grant table [Database_name] [table name] [username] [access type]</code></pre>
    <p>- <code>[Database_name]</code>: The database containing the table.</p>
    <p>- <code>[table name]</code>: The table to grant access to.</p>
    <p>- <code>[username]</code>: The user receiving the access.</p>
    <p>- <code>[access type]</code>: One of <code>[select]</code>, <code>[update]</code>, <code>[insert]</code>, or <code>[delete]</code>.</p>
    <p><strong>Example:</strong></p>
    <pre><code>grant table mydb mytable myuser select</code></pre>

    <h2 id="query">Query</h2>
    <p>After connecting to a database, the <code>query</code> command executes SQL queries. Transaction management is also supported.</p>

    <h3>Execute a Query</h3>
    <p><strong>Command:</strong></p>
    <pre><code>query [SQL_query]</code></pre>
    <p>- <code>[SQL_query]</code>: The SQL query to execute.</p>
    <p><strong>Example:</strong></p>
    <pre><code>query SELECT * FROM mytable</code></pre>

    <h3>Transaction Management</h3>
    <p>Use these commands to manage transactions:</p>
    <ul>
        <li><strong><code>start</code></strong>: Begins a new transaction.
            <pre><code>start</code></pre>
        </li>
        <li><strong><code>commit</code></strong>: Commits the current transaction.
            <pre><code>commit</code></pre>
        </li>
        <li><strong><code>rollback</code></strong>: Rolls back the current transaction.
            <pre><code>rollback</code></pre>
        </li>
    </ul>
    <p><strong>Example with Transactions:</strong></p>
    <pre><code>start
query INSERT INTO mytable (col1, col2) VALUES (1, 'a')
query UPDATE mytable SET col2 = 'b' WHERE col1 = 1
commit</code></pre>

    <h2 id="link-and-migrate">Link and Migrate</h2>
    <p>These are advanced commands that may require additional configuration.</p>

    <h3>Link</h3>
    <p><strong>Command:</strong></p>
    <pre><code>link [Database_name] [postgres connection string]</code></pre>
    <p>- <code>[Database_name]</code>: The DuckDB database to link.</p>
    <p>- <code>[postgres connection string]</code>: A connection string to a PostgreSQL database.</p>
    <p><strong>Purpose:</strong> Likely establishes a connection between DuckDB and PostgreSQL databases.</p>
    <p><strong>Example:</strong></p>
    <pre><code>link mydb "postgresql://user:password@localhost:5432/pgdb"</code></pre>

    <h3>Migrate</h3>
    <p><strong>Command:</strong></p>
    <pre><code>migrate [Database_name]</code></pre>
    <p>- <code>[Database_name]</code>: The database to migrate.</p>
    <p><strong>Purpose:</strong> Possibly migrates data or schema; exact functionality is unspecified.</p>
    <p><strong>Example:</strong></p>
    <pre><code>migrate mydb</code></pre>
    <p><strong>Note:</strong> Consult additional documentation or server logs for clarification on <code>link</code> and <code>migrate</code>.</p>

    <h2 id="usage-example">Usage Example</h2>
    <p>Here’s a step-by-step example of using the server:</p>
    <ol>
        <li><strong>Login as Super User:</strong>
            <pre><code>login duck superpassword</code></pre>
        </li>
        <li><strong>Create a Database:</strong>
            <pre><code>create database mydb</code></pre>
        </li>
        <li><strong>Create a User:</strong>
            <pre><code>create user myuser 12345678</code></pre>
        </li>
        <li><strong>Grant Database Access:</strong>
            <pre><code>grant database mydb myuser read</code></pre>
        </li>
        <li><strong>Login as New User:</strong> (Assuming a new session or reconnection)
            <pre><code>login myuser 12345678</code></pre>
        </li>
        <li><strong>Connect to Database:</strong>
            <pre><code>connect mydb</code></pre>
        </li>
        <li><strong>Execute a Query:</strong>
            <pre><code>query SELECT * FROM mytable</code></pre>
        </li>
    </ol>

    <h2 id="security-and-access-control">Security and Access Control</h2>
    <p>The server enforces strict access control:</p>
    <ul>
        <li>Only the super user <code>duck</code> can create databases and users or grant permissions.</li>
        <li>Users can only connect to databases and execute queries on tables they have been granted access to.</li>
        <li>Permissions are checked for every connection and query operation.</li>
    </ul>
    <p>This documentation provides the essentials for using the TCP Server for DuckDB. For further details, such as starting the server or handling advanced features, refer to the project’s additional documentation or source code.</p>
</body>
</html>