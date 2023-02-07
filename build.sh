#!/usr/bin/env bash
echo "Chaỵ Binary BSServiceMap trong terminal"

swag swag init --dir ./ --generalInfo clientgin/client1.go --propertyStrategy snakecase --output ./routes/docs;

echo "Build Server"
# shellcheck disable=SC2164
cd server
go build



echo "Build Client"
# shellcheck disable=SC2164
cd ../client
go build

cd ..

