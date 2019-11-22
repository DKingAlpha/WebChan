export GOPATH=`dirname $(realpath $0)`
go build -o $GOPATH/webchan cmd/webchan
