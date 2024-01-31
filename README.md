![Build and Lint Windhoek](https://github.com/runvelocity/windhoek/actions/workflows/build-windhoek.yaml/badge.svg)

<div align="center">
  <!-- <a href="https://github.com/othneildrew/Best-README-Template">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a> -->

  <h3 align="center">Windhoek</h3>

  <p align="center">
    Windhoek is a small API written in Golang that exposes routes to create Firecracker VMs to run velocity functions.
    <br />
    <a href="https://docs.runvelocity.dev"><strong>Explore the docs »</strong></a>
    <br />
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about">About</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

## About

Windhoek is an API used as part of the velocity stack to manage the creation of Firecracker MicroVMs and running Velocity functions within them

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

This section outlines how to run the Windhoek API within your Virtual machine.

### Prerequisites

This API creates Firecracker MicroVMs and these VMs can only be created on Virtual machines that support nested virtualization, such as, Digitalocean droplets, AWS .metal instances and bare-metal servers. 


### Installation

A convinience installation script has been provided and you can quicly run the API by using the following command

```bash
curl https://gist.githubusercontent.com/utibeabasi6/6cdf2c52262ef8eeffba2e4d6970d36c/raw/388413f24c06cd79017d948a4f6e049cf556d6f7/windhoek-setup.sh | bash
```

Otherwise, you can install the individual components by following the steps below.

<details>
<summary>Install Golang</summary>
<br/>

```bash
 pushd /tmp
 wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
 rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
 echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
 popd
```
</details>

<details>
<summary>Install Firecracker</summary>
<br/>

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
</details>

<details>
<summary>Install CNI plugins</summary>
<br/>

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
</details>

<details>
<summary>Create required directories</summary>
<br/>


```bash
mkdir /root/fckernels /root/fcsockets /root/fcruntimes
```

</details>

<details>
<summary>Download Kernel and Runtime root filesystem</summary>
<br/>

```bash
cd /root/fcruntimes && wget https://terraform-20231223074656017300000001.s3.amazonaws.com/runtimes/nodejs-runtime.ext4
cd /root/fckernels && wget https://terraform-20231223074656017300000001.s3.amazonaws.com/kernels/vmlinux

```

</details>

<details>
<summary>Setup systemd</summary>
<br/>

#### Create the systemd service 

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
</details>

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Usage

### Invoking the API

The Windhoek API is started on port 8000. To run a function, send a POST request to `http://<your-ip>:8000/invoke` with the following payload

```json
{
    "args": {}, // Arguments to the function
    "codeLocation": "", // S3 URL where the function code is stored
    "handler": "" // The function's entrypoint (filename without extension)
}
```

_For more examples, please refer to the [Documentation](https://docs.runvelocity.dev)_

<!-- ROADMAP -->
## Roadmap

- [x] Add Function invocations
- [ ] Add endpoints to delete MicroVms
- [ ] Keep execution environments warm

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->
## Contact

Utibeabasi Umanah - [@utibeumanah_](https://twitter.com/utibeumanah_) - utibeabasiumanah6@gmail.com

Project Link: [https://github.com/runvelocity](https://github.com/runvelocity)

<p align="right">(<a href="#readme-top">back to top</a>)</p>