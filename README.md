# Tcp Server For DuckDB

## Communication protocol
> [!Caution]
>  ### **all commend needs to be sent in lower case.**

### `connect` $\rightarrow$ `$dbname$`
- ####  using connect requires first to login.
- #### you can not connect to a database which you are not authorized to access within the server.
- #### if connected successfully the response would be `"success"`  