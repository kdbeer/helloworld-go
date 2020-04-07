package main

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf just config
var Conf Config

// Config just config
type Config struct {
	ProjectName string
	LogFolder   string
	Mongo       struct {
		Address  string
		Timeout  time.Duration
		Database string
		Username string
		Password string
	}
	Directory struct {
		Root   string
		Images string
		Models string
	}
}

func (c *Config) initViper(stage string) {
	viper.AddConfigPath("./")
	viper.SetConfigFile("config.yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Read config file error: %v", err)
	}

	c.binding(stage)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("%s has changed", e.Name)
		c.binding(stage)
	})
}

func (c *Config) binding(stage string) {
	sub := viper.Sub(stage)
	c.LogFolder = sub.GetString("log_folder")

	c.Mongo.Address = fmt.Sprintf("%s:%v", sub.GetString("mongo.address"), sub.GetInt("mongo.port"))
	c.Mongo.Timeout = sub.GetDuration("mongo.timeout")
	c.Mongo.Username = sub.GetString("mongo.username")
	c.Mongo.Password = sub.GetString("mongo.password")
	c.Mongo.Database = sub.GetString("mongo.schema.salon")

	c.Directory.Root = sub.GetString("directory.root")
	c.Directory.Images = path.Join(c.Directory.Root, "images")
	c.Directory.Models = path.Join(c.Directory.Root, "models")
}
