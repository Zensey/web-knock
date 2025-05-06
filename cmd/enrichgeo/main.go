package main

import (
	repo_ "github.com/zensey/web-knock/repo"
)

func main() {
	repo := new(repo_.Repo)
	repo.Init()
	georepo := new(repo_.Georepo)
	georepo.Init()

	list := repo.GetBlacklist()
	for _, r := range list {
		repo.BlacklistUpdateGeo(r.IP, georepo.Lookup(r.IP))
	}
}
