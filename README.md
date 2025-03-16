# Tcp Server For DuckDB

## Communication protocol
> [!Caution]
>  ### **all commend and names passed are not case sensetive**

### `connect`
<img src="image/connect.svg">

- ##### requires first to login.
- ##### you can not connect to a database which you are not authorized to access within the server.
- ##### if connected successfully the response would be `"success"`
- ##### this is the first step before executing queries and using the datbase.
> [!TIP]
> ```bash
>   connect mydb
> ```

### `create`
<img src="image/create.svg">

- ##### requires first to login only as the super user `duck`.
- ##### `database_name` and `username` are unique per server.
- ##### creating a user does not gave him any access over any database in the server you need to grant him access using `grant` command