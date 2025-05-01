//go:build linux
// +build linux

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/zensey/web-knock/model"

	"github.com/lrh3321/ipset-go"
	strftime "github.com/ncruces/go-strftime"
	"github.com/nxadm/tail"
	"github.com/satyrius/gonx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	setWhite = "ssh_whitelist"
	setBlack = "web_blacklist"

	fileNameDefault = "/var/log/nginx/access.log"

	// access.log sample
	// 164.52.24.188 - - [29/Mar/2025:17:17:55 -0400] "GET /v1/models HTTP/1.1" 404 555 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	logFormat = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\""
)

func ipsetAdd(setname, strIP string) {
	log.Println("ipsetAdd", setname, strIP)

	err := ipset.Create(setname, ipset.TypeHashIP, ipset.CreateOptions{})
	if err != nil && err.Error() != "file exists" {
		log.Println(err)
	}
	ip := net.ParseIP(strIP)
	err = ipset.Add(setname, &ipset.Entry{IP: ip})
	if err != nil && err.Error() != "exist" {
		log.Println(err)
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

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	db.AutoMigrate(&model.Blacklist{})
	p := gonx.NewParser(logFormat)

	map1 := make(map[string]window)

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
			ipsetAdd(setWhite, remote)
			return
		}

		if status == "404" || status == "400" {
			count := putEvent(map1, remote, time)
			if count == 4 {

				rec := &model.Blacklist{IP: remote, Request: time}
				db.Save(rec)
				ipsetAdd(setBlack, remote)
			}
		}
	}

	for {
		func() {
			file, err := os.Open(*fileName)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				t := scanner.Text()
				handleLine(t)
			}
		}()

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
