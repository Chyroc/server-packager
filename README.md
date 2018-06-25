# server-packager
static server that packages all resources together / 将所有资源打包到一起的静态服务器

Serve embedded files from [jteeuwen/go-bindata](https://github.com/jteeuwen/go-bindata) with `http server`.

### Installation

Install with

    $ go get github.com/jteeuwen/go-bindata/...
    $ go get github.com/Chyroc/server-packager

### Creating embedded data

Usage is identical to [jteeuwen/go-bindata](https://github.com/jteeuwen/go-bindata) usage,
instead of running `go-bindata` run `server-packager`.

The tool will create a `bindata.go` file, which contains the embedded data, and `main.go` file which contains http server.

A typical use case is

    $ server-packager data/...
