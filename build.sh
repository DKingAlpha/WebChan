export GOPATH=`dirname $(realpath $0)`
cd $GOPATH
git pull
go build -o $GOPATH/webchan cmd/webchan
