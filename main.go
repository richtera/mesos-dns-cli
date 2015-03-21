package main

import (
	"flag"
	"fmt"
	"github.com/benschw/dns-clb-go/clb"
	"github.com/benschw/dns-clb-go/dns"
	"os"
	"path/filepath"
	"strings"
)

var version = "1.0"

func main() {
	versionFlag := false

	dnsSrv := flag.String("dns", "127.0.0.1:53", "dns to query")
	index := flag.String("index", "", "Index name to add to primary name (hack)")
	protocol := flag.String("protocol", "tcp", "Protocol tcp/udp and so on.")
	framework := flag.String("framework", "marathon", "Framework")
	domain := flag.String("domain", "mesos", "Domain")
	flag.BoolVar(&versionFlag, "version", false, "output the version")
	hname := flag.String("hname", "", "hostname prefix if needed")
	pname := flag.String("pname", ".", "port prefix if needed")
	flag.Parse()

	if versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}
	if flag.NArg() != 1 {
		os.Exit(1)
	}

	srvName := flag.Args()[0]
	firstChar := string(srvName[0])
	realQuery := srvName
	if firstChar == "." || firstChar == "/" {
		pathInfo := ""
		if firstChar == "." {
			pathInfo = filepath.Clean(filepath.Join(os.Getenv("MARATHON_APP_ID"), srvName))
		} else {
			pathInfo = realQuery
		}
		source := strings.Split(pathInfo, "/")
		query := ""
		for i := len(source) - 1; i >= 0; i-- {
			s := source[i]
			if s != "" {
				if query != "" {
					query = query + "." + s
				} else {
					query = "_" + s
					if index != nil && *index != "" {
						query = query + "_" + *index
					}
				}
			}
		}
		realQuery = query + "._" + *protocol + "." + *framework + "." + *domain + "."
	}
	lib := dns.NewLookupLib(*dnsSrv)
	c := clb.NewRandomClb(lib)
	address, err := c.GetAddress(realQuery)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	if *hname != "" {
		fmt.Printf("%s%s%s%d", *hname, address.Address, *pname, address.Port)
	} else {
		fmt.Printf("%s", address)
	}
}
