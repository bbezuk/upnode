package upnode

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	"log"
)

func CreateTableRequest(db *sql.DB) (err error) {
	_, err = db.Exec(` CREATE TABLE "Request"
			 (
	 		 "requestId" smallint NOT NULL,
	 		 "interval" smallint,
	 		 type integer NOT NULL,
	 		 address character varying(200) NOT NULL,
	 		 "expectedResult" text,
	 		 options character varying(100),
	 		 priority smallint NOT NULL,
			 push smallint NOT NULL DEFAULT 0,
	 		 one_time smallint NOT NULL DEFAULT 0,
	 		 paused smallint NOT NULL DEFAULT 0,
	 		 seq integer NOT NULL DEFAULT 0,
	 		 "backupSeq" integer NOT NULL DEFAULT 0,
	 		 "plantedSeed" smallint,
	 		 "fallbackList" smallint[],
	 		 "slotList" character varying(8)[],
	 		 updated integer NOT NULL DEFAULT 0,
	 		 CONSTRAINT "Request_pkey" PRIMARY KEY ("requestId" )
		 )`)

	if err != nil {
		log.Printf("'Request' table create query failed: \n%s\n", err)
	}

	return
}

func DropTableRequest(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE \"Request\" ")

	if err != nil {
		log.Printf("'Request' table drop query failed: \n%s\n", err)
	}

	return
}

func CreateTableCluster(db *sql.DB) (err error) {
	_, err = db.Exec(`CREATE TABLE "Cluster"
		 (
		  id serial NOT NULL,
		  ip character varying(20),
		  port integer NOT NULL DEFAULT 13339,
		  "inputCode" character varying(20),
		  "outputCode" character varying(20),
		  key character varying(128),
		  active smallint NOT NULL DEFAULT 0,
		  "group" smallint,
		  CONSTRAINT "Cluster_pkey" PRIMARY KEY (id )
	  ) `)

	if err != nil {
		log.Printf("Table create query failed: \n%s\n", err)
	}

	return
}

func DropTableCluster(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE \"Cluster\" ")

	if err != nil {
		log.Printf("Table drop query failed: \n%s\n", err)
	}

	return
}

func CreateTableIncident(db *sql.DB) (err error) {
	_, err = db.Exec(`CREATE TABLE "Interval"
		 (
		  id serial NOT NULL,
		  value integer,
		  CONSTRAINT "Interval_pkey" PRIMARY KEY (id )
	  ) `)

	if err != nil {
		log.Printf("Table create query failed: \n%s\n", err)
	}

	return
}

func DropTableIncident(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE \"Interval\" ")

	if err != nil {
		log.Printf("Table drop query failed: \n%s\n", err)
	}

	return
}

func InsertDefaultIntervals(db *sql.DB) (err error){
	_, err = db.Exec(`INSERT INTO "Interval" VALUES 
		(1,1),
		(2,2),
		(3,5),
		(4,10),
		(5,20),
		(6,30),
		(7,60)
	`)

	if err != nil {
		log.Printf("Default intervals fill failed: \n%s\n", err)
	}

	return
}
