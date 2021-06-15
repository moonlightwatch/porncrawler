package data

import (
	"context"
	"encoding/json"

	"time"

	"github.com/go-redis/redis/v8"
)

func NewDataInterface(opt *redis.Options) *DataInterface {
	d := &DataInterface{}
	d.QueueDB = redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       0,
	})
	d.PornDB = redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       1,
	})
	return d
}

type DataInterface struct {
	QueueDB *redis.Client
	PornDB  *redis.Client
}

func (d *DataInterface) Close() {
	d.PornDB.Close()
	d.QueueDB.Close()
}
func (d *DataInterface) GetTarget() string {
	ctx := context.TODO()
	key, err := d.QueueDB.RandomKey(ctx).Result()
	if err != nil {
		return ""
	}
	url, err := d.QueueDB.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	d.QueueDB.Del(ctx, key).Result()
	return url
}

func (d *DataInterface) AddTarget(domain, url string) {
	ctx := context.TODO()
	dur, _ := time.ParseDuration("0s")
	d.QueueDB.Set(ctx, domain, url, dur).Result()
}

func (d *DataInterface) AddSite(s SiteData) {
	ctx := context.TODO()
	siteData, err := json.Marshal(map[string]interface{}{
		"title": s.Title,
		"url":   s.URL,
		"time":  time.Now().Unix(),
	})
	if err != nil {
		return
	}
	d.PornDB.Set(ctx, s.URL, siteData, 240*time.Hour).Result() // 存十天
}
