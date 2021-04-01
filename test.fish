#!/usr/bin/env fish

echo "Building..."
go build

echo "Running test cases..."
for file in (ls -p testcases | grep -v '/$')
    echo ">> ./testcases/$file"
    ./nuledger < "./testcases/$file" | diff - "./testcases/expouted/$file"
end
