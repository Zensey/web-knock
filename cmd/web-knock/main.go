package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/zensey/web-knock/model"
	"github.com/zensey/web-knock/repo"
	"github.com/zensey/web-knock/utils"

	strftime "github.com/ncruces/go-strftime"
	"github.com/nxadm/tail"
	"github.com/satyrius/gonx"
)

const (
	setWhite = "ssh_whitelist"
	setBlack = "web_blacklist"
	banTTL   = time.Hour

	fileNameDefault = "/var/log/nginx/access.log"

	// access.log sample
	// 164.52.24.188 - - [29/Mar/2025:17:17:55 -0400] "GET /v1/models HTTP/1.1" 404 555 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	logFormat = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\""
)

func blacklistRecreateFromIpset(repo *repo.Repo) {
	s, _ := utils.IpsetGet(setBlack)
	for _, r := range s {
		log.Println(r.IP)
		repo.BlacklistCreateIfNE(r.IP.String())
	}
}

func blacklistClearExpiredBans(repo *repo.Repo) {
	for {
		list := repo.GetBlacklist()
		for _, r := range list {
			if time.Since(r.Request) > banTTL {
				utils.IpsetDel(setBlack, r.IP)
				// repo.blacklistDel(r.IP)
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func main() {
	fileName := flag.String("f", fileNameDefault, "file path")
	key := flag.String("key", "", "auth key")
	flag.Parse()
	if *key == "" {
		log.Println("key must be set")
		return
	}

	repo := new(repo.Repo)
	repo.Init()
	// blacklistRecreateFromIpset()
	go blacklistClearExpiredBans(repo)

	p := gonx.NewParser(logFormat)
	ipmap := make(map[string]model.Window)

	handleLine := func(l string) {
		e, _ := p.ParseString(l)
		req, _ := e.Field("request") // the request: GET /path HTTP/1.1
		remote, _ := e.Field("remote_addr")
		timeStr, _ := e.Field("time_local")
		time, _ := strftime.Parse(`%d/%b/%Y:%H:%M:%S %z`, timeStr)
		status, _ := e.Field("status")

		reqa := strings.Split(req, " ")
		if len(reqa) == 3 && reqa[1] == *key {
			// Knock
			utils.IpsetAdd(setWhite, remote)
			return
		}

		if strings.HasPrefix(status, "40") { // 40x
			count := model.PutEvent(ipmap, remote, time)

			if count == 4 {
				rec := &model.Blacklist{IP: remote, Request: time}
				repo.SaveBlacklist(rec)
				utils.IpsetAdd(setBlack, remote)
			}
		}
	}

	for {
		t, err := tail.TailFile(*fileName, tail.Config{Follow: true, MustExist: false})
		if err != nil {
			log.Fatal(err)
		}
		for line := range t.Lines {
			handleLine(line.Text)
		}

		fmt.Println("wait 1 sec & try to reopen the log")
		time.Sleep(1 * time.Second)
	}
}
