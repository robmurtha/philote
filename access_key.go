package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type AccessKey struct {
	Write []string `json:"write"`
	Read  []string `json:"read"`
	Token string   `json:"-"`
}

func LoadKey(token string) (*AccessKey, error) {
	ak := &AccessKey{Token: token}
	r := RedisPool.Get()

	rawKey, err := redis.String(r.Do("GET", "philote:token:" + token))
	r.Close()
	if err != nil {
		return ak, err
	}
	if rawKey  == "" {
		return ak, InvalidTokenError{"unknown token"}
	}


	err = json.Unmarshal([]byte(rawKey), &ak); if err != nil {
		return ak, InvalidTokenError{"invalid token data: " + err.Error()}
	}

	return ak, nil
}

func (ak *AccessKey) CanWrite(channel string) bool {
	for _, c := range ak.Write {
		if c == "channel" {
			return true
		}
	}
	
	return false
}