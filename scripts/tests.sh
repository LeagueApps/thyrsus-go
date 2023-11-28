#!/usr/bin/env bash -e -i

unset API_BASE_URL
unset INTERNAL_API_BASE_URL
set -e
go test -v -race -cover -coverprofile=profile.cov ./pkg/thyrsus
set +e
echo ""
echo "=================================================================================="
go tool cover -func=profile.cov
echo "=================================================================================="
go tool cover -func=profile.cov | tail -n 1 | awk '{print $3}' | sed s/%// | awk 'BEGIN { FS = "." } ; {print $1}' > cover.num
TOTAL=$(cat cover.num)
THRESHOLD=80
go tool cover -html=profile.cov -o coverage.html
if [[ $TOTAL -ge $THRESHOLD ]]; then
    echo "Passed coverage threshold of ${THRESHOLD}%"
    if [[ $THRESHOLD -lt 80 ]]; then
        echo "That THRESHOLD could be higher, I believe in you! 🫠"
    fi
    exit 0
else
    echo "FAILED coverage threshold of ${THRESHOLD}%"
    if [[ $THRESHOLD -lt 80 ]]; then
        echo "That THRESHOLD sould be higher, it's sad you missed such an easy target 😕"
    fi
    exit 20
fi