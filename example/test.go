package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/halimath/depot"
	"github.com/halimath/depot/example/repo"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// pool, err := sql.Open("mysql", "mixingnotes:password@tcp(localhost:3306)/mixingnotes?parseTime=true")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // See "Important settings" section.
	// pool.SetConnMaxLifetime(time.Minute * 3)
	// pool.SetMaxOpenConns(10)
	// pool.SetMaxIdleConns(10)

	prepareDB()

	factory, err := depot.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer factory.Close()

	repo := &repo.MessageRepo{}

	ctx := context.Background()
	session, ctx, err := factory.Session(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := session.CommitIfNoError(); err != nil {
			log.Fatal(err)
		}
	}()

	msg, err := repo.LoadByID(ctx, "1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", msg)

	msg.Text = "hello, one more time"
	msg.Created = time.Now()

	if err := repo.Update(ctx, msg); err != nil {
		log.Fatal(err)
	}

	msg, err = repo.LoadByID(ctx, "1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", msg)
}

func prepareDB() {
	os.Remove("./test.db")
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("create table messages (id varchar primary key, text varchar, order_index integer, len float, attachment blob, created timestamp)")
	if err != nil {
		log.Fatal(err)
	}

	insert, err := db.Prepare("insert into messages values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insert.Close()

	_, err = insert.Exec("1", "hello, world", 1, 1.1, []byte{1, 2, 3}, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	_, err = insert.Exec("2", "hello, again", 2, 1.1, []byte{1, 2, 3}, time.Now())
	if err != nil {
		log.Fatal(err)
	}
}
