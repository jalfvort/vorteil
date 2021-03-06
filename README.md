<br />
<p align="center">
  <a href="https://github.com/vorteil/vorteil">
    <img src="assets/images/vlogo.png" alt="vorteil">
  </a>
  <h3 align="center">vorteil.io</h3>
    <h5 align="center">turn your applications and containers into micro virtual machines</h5>
</p>
<hr/>

[![Build Status](https://travis-ci.org/vorteil/vorteil.svg?branch=master)](https://travis-ci.org/vorteil/vorteil) <a href="https://codeclimate.com/github/vorteil/vorteil/maintainability"><img src="https://api.codeclimate.com/v1/badges/bb819f04d7550b91f24d/maintainability" /></a> [![Go Report Card](https://goreportcard.com/badge/github.com/vorteil/vorteil)](https://goreportcard.com/report/github.com/vorteil/vorteil) [![](https://godoc.org/github.com/vorteil/vorteil?status.svg)](http://godoc.org/github.com/vorteil/vorteil) [![Discord](https://img.shields.io/badge/chat-on%20discord-6A7EC2)](https://discord.gg/VjF6wn4)

Vorteil is an operating system for running cloud applications on micro virtual machines. It takes only the files you need and runs them on a custom Linux kernel without any unnecessary background services: there's no ssh, no shell, and no login; just a toml file that Vorteil's init process (vinitd) uses to configure the system and launch your apps.

<p align="center">
  <a href="https://github.com/vorteil/vorteil">
    <img src="assets/images/terminal.gif" alt="terminal">
  </a>
</p>

The tools that build Vorteil images have been optimized for speed, and can take your project from loose files to a running virtual machine in as little as one second.

## Installation

The Vorteil binary runs anywhere, just download the archive for your system, extract it, and put the executable on your `PATH`.

### Linux
```sh
wget https://github.com/vorteil/vorteil/releases/latest/download/vorteil_linux-x86.tar.gz
tar -xzf vorteil_linux-x86.tar.gz
./vorteil version
```
### Windows

```sh
wget https://github.com/vorteil/vorteil/releases/latest/download/vorteil_windows-x86.zip -UseBasicParsing -OutFile .\vorteil_windows-x86.zip
Expand-Archive -Path .\vorteil_windows-x86.zip -DestinationPath .\
.\vorteil.exe version
```

### Mac

```sh
curl -LJO https://github.com/vorteil/vorteil/releases/latest/download/vorteil_darwin-x86.dmg
hdiutil attach vorteil_darwin-x86.dmg
cp /Volumes/Vorteil\ CLI/vorteil ./vorteil
hdiutil detach /Volumes/Vorteil\ CLI
./vorteil version
```

### Dependencies

To use the `vorteil run` command you'll need [VirtualBox](https://www.virtualbox.org/wiki/Downloads), [QEMU](https://www.qemu.org/download/), [firecracker](https://github.com/firecracker-microvm/firecracker) or Hyper-V installed on your system and reachable on the `PATH`.

If you're using Windows, it's recommended that you enable developer mode as well, so that the tools can use Unix-style symlinks.

## Using Vorteil

### Hello World

```sh
vorteil run https://apps.vorteil.io/file/vorteil/helloworld
```

This command downloads our Hello World package from apps.vorteil.io and runs it using one of your local hypervisors.

Our Hello World app is a simple web server that should be reachable via NAT on port 8888 (http://localhost:8888/). If that port was unavailable you should notice a yellow warning message in the logs telling you which port it bound instead.

### Convert from a Container

An easy way to build Vorteil micro virtual machines is to build them from container images. Vorteil can convert all OCI compliant images from remote repositories or local docker and containerd runtimes. The following command would convert the _hello-world_ application from docker hub.

```sh

# converts hello-world from docker hub into /tmp/hello
vorteil projects convert-container hello-world /tmp/hello

# converts hello-world from local docker into /tmp/hellolocal
vorteil projects convert-container local.docker/hello-world /tmp/hellolocal

# run it
vorteil run /tmp/hellolocal

```

### Deploying an Image

Vorteil can deploy the images to different cloud providers. Please read the [documentation](https://docs.vorteil.io/docs/guides/configure-provisioner/) for examples.


### Modify an Existing Package

```sh
wget -O cockroachdb.vorteil https://apps.vorteil.io/file/vorteil/cockroachdb
vorteil unpack cockroachdb.vorteil cockroachdb
ls cockroachdb
```

These commands download a Vorteil package and use the `vorteil unpack` command to extract its contents into a structured Vorteil project directory. Think of this directory like it will become the root directory on your app's virtual machine. You can modify the app's file-system contents by changing the files in this directory.

There are two types of special files in the project directory: `.vorteilproject` and files that end with the `.vcfg` file extension. These are used by Vorteil to understand how to build the virtual machine and what it should do when it's running.

Make whatever changes you need to make, then repackage them.

```sh
vorteil pack cockroachdb -o cockroachdb-modified.vorteil
```

This `vorteil pack` command packs up your project directory into a package file again.

```sh
vorteil run cockroachdb-modified.vorteil
```

Try out your modified package by using the `vorteil run` command on it.

## Building From Source

These tools are 100% written in Go, which means compiling them is the same as compiling most simple Go programs.

```sh
go build github.com/vorteil/vorteil/cmd/vorteil
```

### Code of Conduct

We have adopted the [Contributor Covenant](https://github.com/vorteil/.github/blob/master/CODE_OF_CONDUCT.md) code of conduct.

### Contributing

Any feedback and contributions are welcome. Read our [contributing guidelines](https://github.com/vorteil/.github/blob/master/CONTRIBUTING.md) for details.

## License

Distributed under the Apache 2.0 License. See `LICENSE` for more information.

## See Also

* The [Vorteil.io](https://vorteil.io/) website.
* The Vorteil [documentation](https://docs.vorteil.io/).
* The free Vorteil [apps repository](http://apps.vorteil.io/).
* The Vorteil [blog](https://blog.vorteil.io/).
* The [Godoc](https://godoc.org/github.com/vorteil/vorteil) library documentation.
