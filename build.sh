export GOPATH=`dirname $(realpath $0)`
cd $GOPATH
git pull
go build -o $GOPATH/webchan cmd/webchan
sudo setcap 'cap_net_bind_service=+ep' $GOPATH/webchan
