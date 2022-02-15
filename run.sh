set -ex
go test -v
exe=common_line
go build -o $GOPATH/bin/$exe
$exe -file="$FILE" -count="${COUNT-5}"