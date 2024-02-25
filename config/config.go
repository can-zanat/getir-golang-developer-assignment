package config

import "github.com/k0kubun/pp/v3"

type Config struct {
	Mongo      Mongo
	ServerPort string
}

func New() (*Config, error) {
	config := &Config{
		Mongo: Mongo{
			URI:      "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true",
			Database: "getircase-study",
		}}

	return config, nil
}

func (c *Config) Print() {
	_, _ = pp.Println(c)
}
