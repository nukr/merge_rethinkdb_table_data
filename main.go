package main

import (
	"io/ioutil"
	"log"

	"encoding/json"

	r "gopkg.in/dancannon/gorethink.v2"
)

func main() {
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatal(err)
	}
	cursor, err := r.DB("2910eb12_d64a_49cc_b2be_54201441e27b").TableList().Run(session)
	if err != nil {
		log.Fatal(err)
	}

	var tableName string
	var tableList []string
	for cursor.Next(&tableName) {
		tableList = append(tableList, tableName)
	}

	for _, tableName := range tableList {
		cursor, err := r.DB("2910eb12_d64a_49cc_b2be_54201441e27b").Table(tableName).Reduce(func(left, right r.Term) r.Term {
			return left.Merge(right)
		}).Default(struct{}{}).Run(session)
		if err != nil {
			log.Fatal(err)
		}
		var value interface{}
		cursor.One(&value)
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			log.Fatal(err)
		}
		errWriteFile := ioutil.WriteFile("./tables/"+tableName+".json", jsonBytes, 0644)
		if errWriteFile != nil {
			log.Fatal(errWriteFile)
		}
	}
}
