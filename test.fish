#!/usr/bin/env fish

set program $argv[1]

echo "Running test cases..."
for file in (ls -p testcases | grep -v '/$')
    echo ">> ./testcases/$file"
    $program < "./testcases/$file" | diff - "./testcases/expouted/$file"
end
