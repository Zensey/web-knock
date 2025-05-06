package utils

import (
	"log"
	"net"

	"github.com/lrh3321/ipset-go"
)

func IpsetAdd(setname, strIP string) error {
	log.Println("+", setname, strIP)

	err := ipset.Create(setname, ipset.TypeHashIP, ipset.CreateOptions{})
	if err != nil && err.Error() != "file exists" {
		log.Println(err)
		return err
	}
	ip := net.ParseIP(strIP)
	err = ipset.Add(setname, &ipset.Entry{IP: ip})
	if err != nil && err.Error() != "exist" {
		log.Println(err)
	}
	return err
}

func IpsetGet(setname string) (entries []ipset.Entry, err error) {
	sets, err := ipset.List(setname)
	if err != nil {
		return nil, err
	}
	return sets.Entries, nil
}

func IpsetDel(setname, strIP string) error {
	// log.Println("-", setname, strIP)

	ip := net.ParseIP(strIP)
	err := ipset.Del(setname, &ipset.Entry{IP: ip})
	if err != nil && err.Error() != "exist" {
		log.Println(err)
	}
	return err
}
