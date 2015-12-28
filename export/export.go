package export

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"passwd/config"
	"passwd/impt"
	"path/filepath"
	"strconv"
)

const EXPORT_FILE_NAME = "exported_password.txt"

func SingleThreadExportPassword() {
	cfg := config.ReadConfig()
	url := cfg.Mg.Host + ":" + strconv.Itoa(int(cfg.Mg.Port))
	database := cfg.Mg.Database
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB with %s, err: %s", url, err)
	}
	defer session.Close()

	db := session.DB(database)
	cns, err := db.CollectionNames()
	if err != nil {
		log.Fatalf("Failed to get collection names, err: %s", err)
	}

	//save all the unique password
	var allPassword map[string]bool = make(map[string]bool, 40000000)
	var account impt.Account
	for _, cn := range cns {
		c := db.C(cn)

		if c.Name == "system.indexes" {
			continue
		}

		log.Println("exporting ", c.Name)
		q := c.Find(nil).Select(bson.M{"password": true})
		iter := q.Iter()

		for iter.Next(&account) {
			if !allPassword[account.Password] {
				allPassword[account.Password] = true
			}
		}
		if err := iter.Close(); err != nil {
			log.Fatalf("Failed to iterate %s, err: %s", c.Name, err)
		}
	} // end of all collections

	writeToFile(allPassword)
}

func MultipleThreadExportPassword() {

}

func writeToFile(passwords map[string]bool) {
	file, err := os.Create(EXPORT_FILE_NAME)

	if nil != err {
		fileName, _ := filepath.Abs(EXPORT_FILE_NAME)
		log.Fatalf("Failed to create file: %s with err: %s", fileName, err)
	}
	defer file.Close()

	for line, _ := range passwords {
		if _, err := file.WriteString(line + "\n"); nil != err {
			log.Fatalln("Failed to write to file!")
		}
	}
}
