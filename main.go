package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	_ "github.com/elazarl/go-bindata-assetfs"
	_ "github.com/jteeuwen/go-bindata"
)

func isDebug(args []string) bool {
	flagset := flag.NewFlagSet("", flag.ContinueOnError)
	debug := flagset.Bool("debug", false, "")
	debugArgs := make([]string, 0)
	for _, arg := range args {
		if strings.HasPrefix(arg, "-debug") {
			debugArgs = append(debugArgs, arg)
		}
	}
	flagset.Parse(debugArgs)
	if debug == nil {
		return false
	}
	return *debug
}

func getBinDataFile() (*os.File, *os.File, []string, error) {
	bindataArgs := make([]string, 0)
	outputLoc := "bindata.go"

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-o" {
			outputLoc = os.Args[i+1]
			i++
		} else {
			bindataArgs = append(bindataArgs, os.Args[i])
		}
	}

	tempFile, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return &os.File{}, &os.File{}, nil, err
	}

	outputFile, err := os.Create(outputLoc)
	if err != nil {
		return &os.File{}, &os.File{}, nil, err
	}

	bindataArgs = append([]string{"-o", tempFile.Name()}, bindataArgs...)
	return outputFile, tempFile, bindataArgs, nil
}

func genBindata() {
	path, err := exec.LookPath("go-bindata")
	if err != nil {
		fmt.Println("cannot find go-bindata executable in path")
		fmt.Println("maybe you need: go get https://github.com/jteeuwen/go-bindata/...")
		os.Exit(1)
	}
	out, in, args, err := getBinDataFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot create temporary file", err)
		os.Exit(1)
	}
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error: go-bindata: ", err)
		os.Exit(1)
	}
	r := bufio.NewReader(in)
	for line, isPrefix, err := r.ReadLine(); err == nil; line, isPrefix, err = r.ReadLine() {
		if !isPrefix {
			line = append(line, '\n')
		}
		if _, err := out.Write(line); err != nil {
			fmt.Fprintln(os.Stderr, "cannot write to ", out.Name(), err)
			return
		}
	}

	in.Close()
	out.Close()
	if err := os.Remove(in.Name()); err != nil {
		fmt.Fprintln(os.Stderr, "cannot remove", in.Name(), err)
	}
}

func main() {
	genBindata()

	outputFile, err := os.Create("main.go")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot create file", err)
		os.Exit(1)
	}
	outputFile.WriteString(`package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/elazarl/go-bindata-assetfs"
)

func main() {
	addr := flag.String("a", ":5000", "address to serve(host:port)")
	prefix := flag.String("p", "/", "prefix path under")
	root := flag.String("r", ".", "root path to serve")
	flag.Parse()

	var err error
	*root, err = filepath.Abs(*root)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("serving %s as %s on %s", *root, *prefix, *addr)
	http.Handle(*prefix, http.StripPrefix(*prefix, http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "data"})))

	logger := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		http.DefaultServeMux.ServeHTTP(w, r)
	})
	err = http.ListenAndServe(*addr, logger)
	if err != nil {
		log.Fatalln(err)
	}
}
`)
	outputFile.Close()
}
