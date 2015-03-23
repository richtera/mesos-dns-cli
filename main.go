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
	verbose := false

	dnsSrv := flag.String("dns", "127.0.0.1:53", "dns to query")
	index := flag.String("index", "", "Index name to add to primary name (hack)")
	protocol := flag.String("protocol", "tcp", "Protocol tcp/udp and so on.")
	framework := flag.String("framework", "marathon", "Framework")
	domain := flag.String("domain", "mesos", "Domain")
	flag.BoolVar(&verbose, "verbose", false, "verbose")
	hname := flag.String("hname", "", "hostname prefix if needed")
	pname := flag.String("pname", ".", "port prefix if needed")
	flag.Parse()

	if flag.NArg() != 1 {
		os.Exit(1)
	}

	srvName := flag.Args()[0]
	firstChar := string(srvName[0])
	realQuery := srvName
	if firstChar == "." || firstChar == "/" {
		if verbose {
			fmt.Printf("Using relative path from %s\n", os.Getenv("MARATHON_APP_ID"))
		}
		pathInfo := ""
		if firstChar == "." {
			pathInfo = filepath.Clean(filepath.Join(os.Getenv("MARATHON_APP_ID")+"/", srvName))
		} else {
			pathInfo = realQuery
		}
		if verbose {
			fmt.Printf("Calculated path to be %s\n", pathInfo)
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
		realQuery = query + "._" + *protocol + "." + *framework + "." + *domain
	}
	if verbose {
		fmt.Printf("Query SRV %s\n", realQuery)
	}
	lib := dns.NewLookupLib(*dnsSrv)
	c := clb.NewRoundRobinClb(lib)
	address, err := c.GetAddress(realQuery)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if *hname != "" {
		fmt.Printf("%s%s%s%d", *hname, address.Address, *pname, address.Port)
	} else {
		fmt.Printf("%s", address)
	}
}
