# Realtime Chat App

## Setting up database locally
- Ensure you have PostgreSQL installed.
- Ensure you have `sqlx-cli` installed or install via ```cargo install sqlx-cli```
- Set the database URL in .env file.
- Create your database schema in `migrations` directory i.e `schema.sql`
- Run the following commands to spin up your local db
  - ```
    sqlx database create
    ```
  - ```
     sqlx migrate add <DESC>
    ```
