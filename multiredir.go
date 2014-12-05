package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/artyom/autoflags"
)

func main() {
	config := struct {
		Timeout time.Duration `flag:"t,connect timeout"`
		Verbose bool          `flag:"v,be verbose on errors"`
		Redirs  redirections  `flag:"r,redirection to serve: bind_host:port/to_host:port (multiple -r params can be used)"`
	}{
		Timeout: 5 * time.Second,
	}
	if err := autoflags.Define(&config); err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	if len(config.Redirs) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	for _, r := range config.Redirs {
		go serve(r, config.Verbose, config.Timeout)
	}
	select {}
}

func serve(r redir, verbose bool, timeout time.Duration) {
	ln, err := net.Listen("tcp", r.From)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn, r.To, verbose, timeout)
	}
}

func handleConn(conn net.Conn, addr string, verbose bool, timeout time.Duration) {
	defer conn.Close()
	remote, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		if verbose {
			log.Print(err)
		}
		return
	}
	defer remote.Close()
	go io.Copy(conn, remote)
	io.Copy(remote, conn)
}

type redir struct {
	From string
	To   string
}

type redirections []redir

func (r *redirections) String() string { return "" }
func (r *redirections) Set(value string) error {
	p := strings.SplitN(value, "/", 2)
	if len(p) != 2 {
		return fmt.Errorf("invalid parameter, should be in form: bind_host:port,to_host:port")
	}
	*r = append(*r, redir{From: p[0], To: p[1]})
	return nil
}
