package config

import (
	"bytes"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/timshannon/bolthold"
)

//Config structure
type Config struct {
	Debug bool

	Timezone, Version, COOKIE, DB, OS,
	Path, Address string

	BoltHold   *bolthold.Store
	Encryption struct {
		Private []byte
		Public  []byte
	}
}

var config Config

//Get ...
func Get() *Config {
	return &config
}

//Init ...
func Init(yamlConfig []byte) {

	viper.SetConfigType("yaml")
	viper.SetDefault("address", "127.0.0.1:8000")

	var err error
	if yamlConfig == nil {
		viper.SetConfigName("config")
		viper.AddConfigPath("./")  // optionally look for config in the working directory
		err = viper.ReadInConfig() // Find and read the config file
	} else {
		err = viper.ReadConfig(bytes.NewBuffer(yamlConfig))
	}

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	config = Config{}
	config.Debug = false
	config.DB = viper.GetString("db")
	config.OS = viper.GetString("os")
	config.Path = viper.GetString("path")

	if config.BoltHold, err = bolthold.Open(
		config.DB, 0666, nil); err != nil {
		log.Fatalf(err.Error())
		return
	}

	config.COOKIE = viper.GetString("cookie")
	config.Address = viper.GetString("address")
	config.Version = viper.GetString("version")
	config.Timezone = viper.GetString("timezone")

	encrptionKeysMap := viper.GetStringMapString("encryption_keys")
	if encrptionKeysMap != nil {
		config.Encryption.Public, err = Asset(encrptionKeysMap["public"])
		if err != nil {
			log.Fatalf("Error reading public key %v", err)
			return
		}

		config.Encryption.Private, err = Asset(encrptionKeysMap["private"])
		if err != nil {
			log.Fatalf("Error reading private key %v", err)
			return
		}
	}
}
