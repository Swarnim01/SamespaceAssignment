package main

import (
	"log"

	"github.com/gocql/gocql"
)

var session *gocql.Session

func initDB() {
	cluster := gocql.NewCluster("172.17.0.2")
	cluster.Keyspace = "todo"
	cluster.Consistency = gocql.Quorum
	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Unable to connect to ScyllaDB:", err)
	}
}

func closeDB() {
	session.Close()
}
