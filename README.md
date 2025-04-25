# gql - TUI database management tool

> [!NOTE]
> This tool does not have to do anything with GraphQL or Graph Query Language \
> I just needed a name and in 5 seconds came up with this and i started to like it \
> before i knew GraphQL or Graph Query Language were a thing

![](.github/assets/TUI_Screenshot1.png)

## Example for connection string / URI
- after an successful connection the uri is saved in `~/.cache/gql/gql_creds` \
    and auto used in the next start of gql

mysql: \
`username:password@(address)/` \
`forgejo:forgejo@(localhost:3306)/`

postgres: \
`postgresql://username:password@address/databasename` \
`postgresql://crolbar:aoeu@localhost:5432/t?sslmode=disable`
- if you have db's that give "don't accept connection" errors, you can add "&connect_timeout=1" to you're connection string to reduce the halt.

## Support
-   [x] MySQL
-   [x] PostgreSQL
-   [ ] SQLite
