package impt

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"os"
	"passwd/config"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ImportAll() {
	metaDatas := ReadImportConfig()
	for _, metaData := range metaDatas {
		log.Println("Importing ", metaData.FileName)
		Import(&metaData)
	}

	log.Println("Successfully imported all! ")
}

func ReadImportConfig() []ImportMetaData {
	absFile, _ := filepath.Abs("import.cfg")
	log.Println("going to open import config file @ ", absFile)
	fp, err := os.Open("import.cfg")
	if nil != err {
		log.Println("Failed to open: ", absFile)
		os.Exit(2)
	}
	scanner := bufio.NewScanner(fp)
	metaDatas := make([]ImportMetaData, 0, 10)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			parts := strings.Split(line, ",")
			emailIndex, err := strconv.Atoi(parts[3])
			if nil != err {
				log.Fatalln("this line has err: ", line)
			}
			pwdIndex, err := strconv.Atoi(parts[4])
			if nil != err {
				log.Fatalln("this line has err: ", line)
			}
			nameIndex, err := strconv.Atoi(parts[5])
			if nil != err {
				log.Fatalln("this line has err: ", line)
			}
			metaDatas = append(metaDatas, ImportMetaData{parts[0], parts[1], parts[2], emailIndex, pwdIndex, nameIndex})
		}
	}
	return metaDatas
}

func Import(metaData *ImportMetaData) {
	if files, err := ListImportFiles(metaData.FileName); nil == err {
		fmt.Println("total files: ", len(files), files)
		var wg sync.WaitGroup
		ch := make(chan []Any, 300000)
		for _, file := range files {
			wg.Add(1)
			fmt.Println("add 1")
			go LoadFileData(file, metaData, ch, &wg)
		}
		go func() {
			fmt.Println("I am waiting")
			wg.Wait()
			fmt.Println("wait done, going to close")
			close(ch)
		}()
		maxProcs := runtime.GOMAXPROCS(0)
		numCPU := runtime.NumCPU()
		fmt.Println("maxProcs: ", maxProcs, ", numCPU: ", numCPU)
		InsertToDb(metaData, ch)
	} else {
		fmt.Println("Failed to import: ", err)
	}

}

func ListImportFiles(fileName string) ([]string, error) {
	fmt.Println("input ", fileName)
	fl, err := os.Open(fileName)
	if nil != err {
		return nil, errors.New("can not open file: " + fileName)
	}
	defer fl.Close()

	stat, err := fl.Stat()
	if stat.IsDir() {
		files, err := ioutil.ReadDir(fileName)
		if nil != err {
			return nil, errors.New("can not open dir: " + fileName)
		}
		subFiles := make([]string, 0)
		for _, v := range files {
			if !v.IsDir() {
				subFiles = append(subFiles, fileName+string(filepath.Separator)+v.Name())
			}
		}
		return subFiles, nil
	} else {
		return []string{fileName}, nil
	}
}

func LoadFileData(absolutFileName string, metaData *ImportMetaData, ch chan []Any, wg *sync.WaitGroup) {
	fmt.Println("loading file: ", absolutFileName)
	file, err := os.Open(absolutFileName)
	if nil != err {
		os.Exit(2) //TODO
	}
	defer file.Close()

	pFunc := mapParseFunc(metaData)
	scanner := bufio.NewScanner(file)
	var line string
	accounts := make([]Any, 0, 300000)
	for scanner.Scan() {
		line = scanner.Text()
		//fmt.Println(line)
		accounts = append(accounts, pFunc(line, metaData))
		if 300000 <= len(accounts) {
			ch <- accounts
			accounts = make([]Any, 0, 300000)
		}
	}

	ch <- accounts
	wg.Done()
	fmt.Println("Done 1 ")
}

func InsertToDb(metaData *ImportMetaData, ch chan []Any) {
	cfg := config.ReadConfig()
	url := cfg.Mg.Host + ":" + strconv.Itoa(int(cfg.Mg.Port))
	database := cfg.Mg.Database
	session, err := mgo.Dial(url)
	if err != nil {
		os.Exit(1)
	}
	defer session.Close()

	c := session.DB(database).C(metaData.Name)
	c.DropCollection() //drop first

	start := time.Now().UnixNano()
	var count = 0
	for accounts := range ch {
		count += len(accounts)
		fmt.Println("Get and insert ", len(accounts))
		c.Insert(covertToInterface(accounts)[:]...)
	}

	fmt.Println("Used : ", (time.Now().UnixNano() - start))
	fmt.Printf("Write to mongoDB(%d) done!\n", count)
}

func covertToInterface(accounts []Any) (records []interface{}) {
	length := len(accounts)
	records = make([]interface{}, length, length)
	for i, _ := range accounts {
		records[i] = accounts[i]
	}
	return
}
