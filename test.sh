#! /usr/bin/env sh

cd backend/

sqlx database create --database-url postgres://myuser:mypass@localhost/test_chat_app

RUST_TEST_THREADS=1 cargo test

echo "removing used dabatase"
sqlx database drop --database-url postgres://myuser:mypass@localhost/test_chat_app -y
