CURDIR=$(dirname $0)/

go build -o $CURDIR../build $CURDIR../cmd/go-open-music

echo
echo "==> Results:"
ls -hl $CURDIR../build/