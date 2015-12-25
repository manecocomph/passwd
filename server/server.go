package server

import (
	"encoding/json"
	"fmt"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"passwd/config"
	"passwd/impt"
	"strconv"
	"sync"
)

const indexHtml = `<!DOCTYPE HTML>
<html>
<head>
    <title>check if your password is leaked</title>
    <meta charset="utf-8" />
</head>
<body>
    <div style="margin: 80px;text-align: center;">
        <input type="text" placeholder="key word" id="kw" autofocus="">
        <button id="SearchBtn">Search</button>
    </div>
    <div id="result" style="width:800px; margin: 0 auto;">
        <span id="statusResult"></span>
        <table id="resultTab" style="margin-top: 13px;">
            <tr style="text-align:left;"><th>Email</th><th>Password</th><th>User Name</th></tr>
            <tbody id="resultTbody">
            </tbody>
        </table>
    </div>
</body>
<script type="text/javascript">
document.getElementById("SearchBtn").addEventListener("click", function(){
var kw = document.getElementById("kw").value
if (null == kw || kw.length < 4) {
    document.getElementById("statusResult").innerText = "Keyword too short!";
    return;
}
document.getElementById("SearchBtn").setAttribute("disabled", "disabled");
document.getElementById("statusResult").innerText = "Searching ...";
var xhttp = new XMLHttpRequest();
xhttp.onreadystatechange = function() {
    if (xhttp.readyState == 4 && xhttp.status == 200) {
        var data = JSON.parse(xhttp.responseText);
        if (0 === data.Status) {
            document.getElementById("statusResult").innerText = "Found:";
            var tbodyHtml = "";

            var accounts = data.Accounts;
            var length = accounts.length;
            for (var i = 0; i < length; i++) {
                tbodyHtml += "<tr><td>" + accounts[i].Email + "</td><td>" + accounts[i].Password + "</td>";
                if (accounts[i].UserName) {
                    tbodyHtml += "<td>" + accounts[i].UserName + "</td>";
                }
                tbodyHtml += "</tr>";
            }
            document.getElementById("resultTbody").innerHTML = tbodyHtml;
        } else {
            document.getElementById("statusResult").innerText = "No result!";
        }
        document.getElementById("SearchBtn").removeAttribute("disabled");
    }
};
xhttp.open("GET", "query?kw=" + kw, true);
xhttp.send();
});
</script>
</html>`

func Start() {
	//key word query
	http.HandleFunc("/query", func(w http.ResponseWriter, req *http.Request) {
		kws := req.URL.Query()["kw"]
		println("before query: ", kws[0])
		result := query(kws[0])
		fmt.Println(result)
		b, err := json.Marshal(result)
		if err != nil {
			io.WriteString(w, "{status:9}")
			return
		}
		fmt.Println(string(b))
		io.WriteString(w, string(b))
	})

	//all others go to index.html
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, indexHtml)
	})
	cfg := config.ReadConfig()
	err := http.ListenAndServe(":"+strconv.Itoa(int(cfg.Serv.Port)), nil)
	if nil != err {
		log.Printf("Can not start server, err: %s", err)
	}
}

type QueryResult struct {
	Status   int
	Accounts []impt.Account
}

func query(kw string) QueryResult {
	//html escape, please
	cfg := config.ReadConfig()
	url := cfg.Mg.Host + ":" + strconv.Itoa(int(cfg.Mg.Port))
	database := cfg.Mg.Database
	session, err := mgo.Dial(url)
	if err != nil {
		return QueryResult{Status: 1}
	}
	defer session.Close()

	db := session.DB(database)
	cns, err := db.CollectionNames()
	if err != nil {
		return QueryResult{Status: 2}
	}

	var allAccount []impt.Account
	var wg sync.WaitGroup
	ch := make(chan []impt.Account, 300)
	for _, cn := range cns {
		c := db.C(cn)

		if c.Name == "system.indexes" {
			continue
		}
		fmt.Println(c.Name)
		/*
				var accounts []impt.Account
				err := c.Find(bson.M{"email": bson.RegEx{".?" + kw + ".?", ""}}).Limit(50).All(&accounts)
				if err != nil {
					return QueryResult{Status: 3}
				}
		        fmt.Println(accounts)*/

		go singleQuery(kw, c, ch, &wg)
		wg.Add(1)
	}
	go func() {
		fmt.Println("I am waiting")
		wg.Wait()
		fmt.Println("wait done, going to close")
		close(ch)
	}()

	for accounts := range ch {
		allAccount = append(allAccount, accounts...)
	}

	return QueryResult{Status: 0, Accounts: allAccount}
}

func singleQuery(kw string, c *mgo.Collection, ch chan []impt.Account, wg *sync.WaitGroup) {
	var accounts []impt.Account
	err := c.Find(bson.M{"email": bson.RegEx{".?" + kw + ".?", ""}}).Limit(50).All(&accounts)
	if err == nil {
		ch <- accounts
	}
	wg.Done()
}
