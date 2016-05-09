#cp bbbmanager /etc/nginx/sites-enabled
#service nginx restart
#apt-get install -y mongodb xmlstarlet
#cd /srv
#git clone blablabla
#cd 4BBB
#wget https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
#tar -C /usr/local -xzf go1.6.2.linux-amd64.tar.gz
# echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/srv/4BBB/server
cd /srv/4BBB/server
#go get
go build
#./server
