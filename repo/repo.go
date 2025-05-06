package repo

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/zensey/web-knock/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repo struct {
	db *gorm.DB
	mu sync.Mutex
}

func (r *Repo) Init() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	db.AutoMigrate(&model.Blacklist{})
	r.db = db
}

func (r *Repo) SaveBlacklist(rec *model.Blacklist) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.db.Save(rec)
}

func (r *Repo) GetBlacklist() []model.Blacklist {
	r.mu.Lock()
	defer r.mu.Unlock()

	var res []model.Blacklist
	r.db.Find(&res)
	return res
}

func (r *Repo) BlacklistCreateIfNE(ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.db.First(&model.Blacklist{IP: ip}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		rec := &model.Blacklist{IP: ip, Request: time.Unix(0, 0)}
		r.db.Save(rec)
	}
}

func (r *Repo) BlacklistUpdateGeo(ip, geo string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rec := &model.Blacklist{IP: ip}
	if err := r.db.First(rec).Error; err == nil {
		rec.Geo = geo
		r.db.Save(rec)
	}
}

func (r *Repo) BlacklistDel(ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.db.Delete(&model.Blacklist{IP: ip}, ip)
}
