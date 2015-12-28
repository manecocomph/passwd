package main

import (
	"flag"
	"fmt"
	"log"
	"passwd/config"
	"passwd/export"
	"passwd/impt"
	serv "passwd/server"
)

func main() {
	flag.Bool("s", true, "start the server")
	flag.Bool("i", false, "import files to MongoDB")
	flag.Bool("t", false, "test if configuration is fine")
	exportFlag := flag.String("x", "", "export use single | multiple thread")
	flag.Parse()
	var serverPresent, isImptPresent, isTstPresent, isExportPresent bool

	flag.Visit(func(fp *flag.Flag) {
		//fmt.Printf("%s: %s  \t default: %s\n", fp.Name, fp.Value, fp.DefValue)
		if "s" == fp.Name {
			serverPresent = true
		}
		if "i" == fp.Name {
			isImptPresent = true
		}

		if "t" == fp.Name {
			isTstPresent = true
		}

		if "x" == fp.Name {
			isExportPresent = true
		}
	})

	if isTstPresent {
		log.Println("testing configuration")
		config.TestConfig()
		return
	}

	if serverPresent {
		log.Println("starting server") //port?
		serv.Start()
		return
	}

	if isImptPresent {
		log.Println("starting import")
		impt.ImportAll()
		return
	}

	if isExportPresent {
		if "single" == *exportFlag {
			log.Println("starting export with single thread")
			export.SingleThreadExportPassword()
			return
		} else if "multiple" == *exportFlag {
			log.Println("starting export with multiple thread")
			return
		}
	}

	fmt.Println("Usage:")
	fmt.Println("\t -s  start the web server")
	fmt.Println("\t -t  test the config.json")
	fmt.Println("\t -i  import the files")
	fmt.Println("\t -x  export password  -x=single|multiple")
	/*

		knownImportMetaData := make([]impt.ImportMetaData, 3, 3)
		knownImportMetaData = append(knownImportMetaData, impt.ImportMetaData{"renren", "/home/xiatian/下载/passwd/renren.txt", impt.SEP_WHITE_SPACE, 0, 1, -1})
		knownImportMetaData = append(knownImportMetaData, impt.ImportMetaData{"csdn", "/home/xiatian/下载/passwd/csdn.txt", " # ", 2, 1, 0})
		knownImportMetaData = append(knownImportMetaData, impt.ImportMetaData{"duduniu", "/home/xiatian/下载/passwd/duduniu", impt.SEP_WHITE_SPACE, 1, 2, 0})
		knownImportMetaData = append(knownImportMetaData, impt.ImportMetaData{"k7k7k", "/home/xiatian/下载/passwd/7k7k", impt.SEP_WHITE_SPACE, 0, 1, -1})

		impt.Import(&impt.ImportMetaData{"k7k7k", "/home/xiatian/下载/passwd/7k7k", impt.SEP_WHITE_SPACE, 0, 1, -1})
	*/
}
