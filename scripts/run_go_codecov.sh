#!/usr/bin/env bash

FILE_COVERAGE_INFO="go_coverage.txt"
FILE_COVERAGE_HTML="go_coverage.html"

set -e
echo "mode: atomic" > ${FILE_COVERAGE_INFO}

for d in $(go list ./internal... | grep -v vendor); do
    go test -race -coverpkg=./... -coverprofile=profile.out -covermode=atomic "$d"
    if [ -f profile.out ]; then
        sed '1d' profile.out >> ${FILE_COVERAGE_INFO}
        rm profile.out
    fi
done

go tool cover -html=./${FILE_COVERAGE_INFO} -o ./${FILE_COVERAGE_HTML}
echo "Export go coverage report to ${FILE_COVERAGE_HTML}"
