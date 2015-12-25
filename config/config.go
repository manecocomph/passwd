package config

import (
	"encoding/json"
	"labix.org/v2/mgo"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

type Mongo struct {
	Host     string
	Port     uint
	Database string
}

type Server struct {
	Port uint
}

type Configuration struct {
	Mg   Mongo
	Serv Server
}

func TestConfig() {
	cfg := ReadConfig()
	testMongo(cfg)
	testServer(cfg)
}

func testMongo(cfg *Configuration) {
	log.Println("Testing MongoDB ...")
	url := cfg.Mg.Host + ":" + strconv.Itoa(int(cfg.Mg.Port))
	database := cfg.Mg.Database
	session, err := mgo.Dial(url)
	if err != nil {
		log.Printf("Can not conect to MongoDB with url: %s, err: %s\n", url, err)
		return
	}
	defer session.Close()

	db := session.DB(database)
	names, err := db.CollectionNames()
	if nil != err {
		log.Printf("Connected, but can't list collections. err: %s\n", err)
		return
	}

	log.Println("Success connected, and current collestions are : ")
	for _, name := range names {
		if "system.indexes" != name {
			log.Println("\t", name)
		}
	}
	log.Println()
}

func testServer(cfg *Configuration) {
	log.Println("Testing web server config ...")
	in, err := net.Listen("tcp", ":"+strconv.Itoa(int(cfg.Serv.Port)))
	if err != nil {
		log.Printf("Port:%d is occupied by other app, err: %s \n", cfg.Serv.Port, err)
		return
	}
	defer in.Close()

	log.Printf("Port:%d is ok to use. \n", cfg.Serv.Port)
}

var cachedCfg *Configuration

func ReadConfig() *Configuration {
	if nil != cachedCfg {
		return cachedCfg
	}

	absPath, _ := filepath.Abs("config.json")
	log.Printf("Reading config file @ %s \n", absPath)
	file, err := os.Open("config.json")
	if nil != err {
		log.Println("config.json not exist, using default values.")
		return New()
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var cfg Configuration
	if err := decoder.Decode(&cfg); nil != err {
		log.Println("Failed to parse config.json, using default values. err: ", err)
		return New()
	}
	cachedCfg = &cfg
	return &cfg
}

func New() *Configuration {
	return &Configuration{Mongo{"127.0.0.1", 27017, "passwd"}, Server{8080}}
}
