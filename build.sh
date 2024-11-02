#! /usr/bin/env sh

rm -r backend/dist

cd frontend && npm run build

cd ../backend && cargo run
