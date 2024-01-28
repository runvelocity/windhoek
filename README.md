# Windhoek

![Build and Lint Windhoek](https://github.com/runvelocity/windhoek/actions/workflows/build-windhoek.yaml/badge.svg)

Windhoek is a small API written in Golang that exposes an invoke route we can call to run a function within a Firecracker VM.

## How to run

### Install dependencies

#### Install Golang
```bash
 pushd /tmp
 wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
 rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
 echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
 popd
```

#### Install Firecracker
```bash
pushd /tmp
ARCH="$(uname -m)"
release_url="https://github.com/firecracker-microvm/firecracker/releases"
latest=$(basename $(curl -fsSLI -o /dev/null -w  %{url_effective} ${release_url}/latest))
curl -L ${release_url}/download/${latest}/firecracker-${latest}-${ARCH}.tgz \
| tar -xz

# Rename the binary to "firecracker"
mv release-${latest}-$(uname -m)/firecracker-${latest}-${ARCH} /usr/local/bin/firecracker
popd
```

#### Install CNI plugins
```bash
apt-get install make
git clone https://github.com/containernetworking/plugins.git /tmp/cni-plugins

# Move to the plugins directory
pushd /tmp/cni-plugins

# Build the CNI tools
./build_linux.sh

mv bin/* /opt/cni/bin

git clone https://github.com/awslabs/tc-redirect-tap
pushd tc-redirect-tap/
make install
mv tc-redirect-tap /opt/cni/bin
popd
popd
```

#### Create required directories

```bash
mkdir /root/fckernels /root/fcsockets /root/fcruntimes
```

#### Download Kernel and Runtime root filesystem

```bash
cd /root/fcruntimes && wget https://terraform-20231223074656017300000001.s3.amazonaws.com/runtimes/nodejs-runtime.ext4
cd /root/fckernels && wget https://terraform-20231223074656017300000001.s3.amazonaws.com/kernels/vmlinux

```

### Deployment

The windhoek service runs as a systemd process within each worker node. To run, follow these steps:

#### Create systemd file

```bash
cat <<EOF > /etc/systemd/system/windhoek.service
[Unit]
Description=Windhoek
After=network.target

[Service]
ExecStart=/root/windhoek
Restart=always

[Install]
WantedBy=default.target
EOF
```

#### Reload the systemd service 
```bash
systemctl daemon-reload
```

#### Enable and start the service
```bash
systemctl enable windhoek.service
systemctl start windhoek.service
```
