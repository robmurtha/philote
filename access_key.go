package main

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

type AccessKey struct {
	Token       string   `json:"-"`
	Read        []string `json:"read"`
	Write       []string `json:"write"`
	AllowedUses int      `json:"allowed_uses"`
	Uses        int      `json:"uses"`
}

func LoadKey(token string) (*AccessKey, error) {
	ak := &AccessKey{Token: token}
	r := RedisPool.Get()

	rawKey, err := redis.String(r.Do("GET", "philote:access_key:" + token))
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
		if c == channel {
			return true
		}
	}
	
	return false
}

func (ak *AccessKey) Save() error {
	r := RedisPool.Get()
	defer r.Close()

	data, err := json.Marshal(ak); if err != nil {
		return err
	}

	_, err =  r.Do("SET", "philote:access_key:" + ak.Token, string(data))
	return err
}

func (ak *AccessKey) Delete() error {
	r := RedisPool.Get()
	defer r.Close()

	_, err := r.Do("DEL", "philote:access_key:" + ak.Token)
	return err
}

func (ak *AccessKey) UsageIsLimited() bool {
	return ak.AllowedUses != 0
}

func (ak *AccessKey) ConsumeUsage() (error) {
	uses, err :=  Lua.ConsumeTokenUsage(ak.Token)
	ak.Uses = uses

	return err
}
