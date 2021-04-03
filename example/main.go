package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/halimath/depot"
	"github.com/halimath/depot/example/models"
	"github.com/halimath/depot/example/repo"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	prepareDB()

	factory, err := depot.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer factory.Close()

	repo := repo.NewMessageRepo(factory)

	ctx := context.Background()
	ctx, err = repo.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := repo.Commit(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	err = repo.Insert(ctx, &models.Message{
		ID:         "1",
		Text:       "hello, world",
		Attachment: []byte{1, 2, 3},
		Created:    time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

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
}
