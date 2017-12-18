# Linuxkit-server

`linuxkit-server` is light wrapper around [moby](https://github.com/moby/tool) to build [linuxkit](https://github.com/linuxkit/linuxkit) operating systems. You can use it to build for example ARM distros in remote server when you cannot do it locally. [EliotOS](https://github.com/ernoaapa/eliot-os) project uses this to create ARM build in CircleCI.

<sub>Built with ❤︎ by [Erno Aapa](https://github.com/ernoaapa) and [contributors](https://github.com/ernoaapa/eliot/contributors)</sub>

## Usage
`linuxkit-server` provides simple HTTP API, where you can POST the linuxkit yaml definition and it will respond with tar file which contains the build output files.


### API
**POST** `/linuxkit/{name}/build/{format}`

Create new build with `{name}` and create `{format}` output. See [linuxkit]() documentation for all format options

### Example
Here's simple example, download `minimal.yml` from linuxkit repository, post it to `linuxkit-server` for building and untar the result to current directory.

```shell
curl https://raw.githubusercontent.com/linuxkit/linuxkit/master/examples/minimal.yml \
  | curl --fail -X POST --data-binary @- http://localhost:8000/linuxkit/example/build/kernel+initrd \
  | tar xvz
```
