set -ex
exe=common_line
go build -o $GOPATH/bin/$exe
$exe -file="$FILE"