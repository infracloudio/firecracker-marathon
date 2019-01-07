#! /bin/bash
set -ex

dpkg -i /mnt/root/linux*.deb

wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
tar -xvf go1.10.3.linux-amd64.tar.gz
mv go /usr/local
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH


echo 'debian-stretch' > /etc/hostname
passwd -d root
mkdir /etc/systemd/system/serial-getty@ttyS0.service.d/
cat <<EOF > /etc/systemd/system/serial-getty@ttyS0.service.d/autologin.conf
[Service]
ExecStart=
ExecStart=-/sbin/agetty --autologin root -o '-p -- \\u' --keep-baud 115200,38400,9600 %I $TERM
EOF

# cat <<EOF > /etc/network/interfaces.d/eth0
# auto eth0
# allow-hotplug eth0
# iface eth0 inet dhcp
# iface eth0 inet6 auto
# 		#ip addr add dev eth0 172.17.100.1/16
# 		#ip route add default via 172.17.0.1
# EOF