//go:build linux
// +build linux

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/lrh3321/ipset-go"
	"github.com/nxadm/tail"
	"github.com/satyrius/gonx"
)

const set = "ssh_whitelist"

func whitelist(strIP string) {
	log.Println("whitelist", strIP)

	err := ipset.Create(set, ipset.TypeHashIP, ipset.CreateOptions{})
	if err != nil {
		log.Println(err)
	}
	ip := net.ParseIP(strIP)
	err = ipset.Add(set, &ipset.Entry{IP: ip})
	if err != nil {
		log.Println(err)
	}
}

func main() {
	key := flag.String("key", "", "auth key")
	flag.Parse()
	if *key == "" {
		log.Println("key must be set")
		return
	}

	for {
		t, err := tail.TailFile("/var/log/nginx/access.log", tail.Config{Follow: true, MustExist: false})
		if err != nil {
			log.Fatal(err)
		}

		// nginx/access.log sample
		// 164.52.24.188 - - [29/Mar/2025:17:17:55 -0400] "GET /v1/models HTTP/1.1" 404 555 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"

		logFormat := "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\""
		p := gonx.NewParser(logFormat)

		for line := range t.Lines {
			e, _ := p.ParseString(line.Text)
			req, _ := e.Field("request") // the request: GET /path HTTP/1.1
			remote, _ := e.Field("remote_addr")

			reqa := strings.Split(req, " ")
			if len(reqa) == 3 && reqa[1] == *key {
				whitelist(remote)
			}
		}
		fmt.Println("wait 1 sec & try to reopen the log")
		time.Sleep(1 * time.Second)
	}
}
