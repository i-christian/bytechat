#! /usr/bin/env sh

rm -r backend/dist

cd frontend && npm run build

cd ../backend && sqlx database create 

cargo run
