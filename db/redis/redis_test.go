package redis

import (
	"config"
	"testing"
)

func TestRedis(t *testing.T) {

	addr, _ := config.Instance().GetValue("test", "redis_addr") //ignore err
	pw, _ := config.Instance().GetValue("test", "redis_passwd") //ignore err
	db, _ := config.Instance().Int("test", "redis_db")          //ignore err
	client, err := NewRedisClient(addr, pw, db)
	if err != nil {
		t.Error("NewRedisClient err, ", err)
	}

	if err = client.Set("t1", "80"); err != nil {
		t.Error("client.Set err, ", err)
	}

	if val, err := client.Get("t1"); err != nil {
		t.Error("client.Get err, ", err)
	} else {
		if val != "80" {
			t.Error("client.Get val is not equal err, ", err)
		}
	}

	if _, err := client.Get("tn"); err != Nil {
		t.Error("client.Get val is not null err, ", err)
	}

	if val, err := client.Exist("t1"); err != nil {
		t.Error("client.Exist err, ", err)
	} else {
		if val != true {
			t.Error("Key Exist!! err, ", err)
		}
	}

	if val, err := client.Exist("tn"); err != nil {
		t.Error("client.Exist err, ", err)
	} else {
		if val != false {
			t.Error("Key not Exist!! err, ", err)
		}
	}

	if err := client.HSet("h1", "f1", "v1"); err != nil {
		t.Error("client.HSet err, ", err)
	}

	if err := client.HSet("h1", "f2", "v2"); err != nil {
		t.Error("client.HSet err, ", err)
	}

	if val, err := client.HGet("h1", "f1"); err != nil {
		t.Error("client.HGet err, ", err)
	} else {
		if val != "v1" {
			t.Error("client.HGet val is not equal err, ", err)
		}
	}

	if _, err := client.HGet("hn", "hn"); err != Nil {
		t.Error("client.HGet val is not null err, ", err)
	}

	if val, err := client.HGetAll("h1"); err != nil {
		t.Error("client.HGet err, ", err)
	} else {
		if len(val) != 2 {
			t.Error("client.HGetAll size is not right err, ", err)
		}
	}

	if val, err := client.HGetAll("hn"); len(val) != 0 {
		t.Error("client.HGetAll val is not null err, ", err)
	}
}
