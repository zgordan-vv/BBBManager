apt-get install -y mongodb xmlstarlet
cd /srv
cd BBBManager/install
cp -f bbbmanager /etc/nginx/sites-enabled
cp -f web /etc/bigbluebutton/nginx
service nginx reload
wget https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.6.2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/root/.gopath
cd ../server
go get
go build
