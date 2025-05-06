package repo

import (
	"log"
	"net/netip"

	"github.com/oschwald/maxminddb-golang/v2"
)

type Georepo struct {
	db *maxminddb.Reader
}

func (r *Georepo) Init() {
	db, err := maxminddb.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	r.db = db
}

type Country struct {
	GeoNameID uint              `maxminddb:"geoname_id"`
	IsoCode   string            `maxminddb:"iso_code"`
	Names     map[string]string `maxminddb:"names"`
}
type City struct {
	Names map[string]string `maxminddb:"names"`
}

type Record struct {
	Country Country `maxminddb:"country"`
	City    City    `maxminddb:"city"`
}

func (r *Georepo) Lookup(ip string) string {
	addr := netip.MustParseAddr(ip)

	var record Record

	err := r.db.Lookup(addr).Decode(&record)
	if err != nil {
		log.Panic(err)
	}
	city := record.City.Names["en"]
	country := record.Country.IsoCode
	if city != "" {
		country = country + ", " + city
	}
	return country
}
