#!/bin/bash

echo "Creating DB."
sleep 5

echo "Building app binary..."
go build -o app .

echo "Runing DB seeder."
./app -seed

echo "Starting DB."
exec ./app

# exec go run .