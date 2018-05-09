echo ">> run golint"
golint

echo ">> run test"
if [ $# -eq 0 ]
then
    runs=1
else
    runs=$1
fi

for ((i=1; i<=$runs; i++))
do
    echo "run #$i"
    go clean -testcache
    go test -cover -race ./smpp
done
