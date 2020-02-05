package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/common-nighthawk/go-figure"
)

var (
	web = false
)

func init() {
	flag.BoolVar(&web, "web", false, "enable webapi")
}

func main() {
	flag.Parse()

	if web {
		server()
	} else {
		panic("Not implimented yet!")
	}
}

func print(w io.Writer, banner, title, ip string) {
	if banner != "" {
		fmt.Fprintf(w, "%s\n", banner)
	}

	if title != "" {
		fmt.Fprintf(w, "%s\n", title)
	}

	if ip != "" {
		fmt.Fprintf(w, "🌐%s\n", ip)
	}
}

func server() {
	sm := http.NewServeMux()
	sm.Handle("/", http.HandlerFunc(handler))
	l, err := net.Listen("tcp4", ":10101")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ready! >>> 0.0.0.0:10101")
	log.Fatal(http.Serve(l, sm))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var banner, title, ip string

	if r.URL.Query().Get("title") != "" {
		title = r.URL.Query().Get("title")
	}

	uriSegments := cleanEmpty(strings.Split(r.URL.Path, "/"))
	if r.URL.Path == "/" || len(uriSegments) == 0 {
		goto End
	} else {
		banner = figure.NewFigure(uriSegments[0], "puffy", false).String()
	}

End:
	if q, _ := strconv.ParseBool(getValueWithDefault(r.URL.Query().Get("ip"), "1")); q {
		ip = getRealIPAdress(r)
	}

	print(w, banner, title, ip)
}

func getRealIPAdress(r *http.Request) string {
	var ipAddress string
	for _, h := range []string{"X-Forwarded-For", "X-Real-IP"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			if ip != "" {
				ipAddress = net.ParseIP(strings.Replace(ip, " ", "", -1)).String()
			}
		}
	}
	if ipAddress == "" {
		ra, _, _ := net.SplitHostPort(r.RemoteAddr)
		ipAddress = ra
	}
	return ipAddress
}

func getValueWithDefault(v string, dv string) string {
	if v == "" {
		return dv
	}
	return v
}

func cleanEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
