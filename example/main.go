// Copyright 2021 Alexander Metzner.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains a cli demonstrating the generated repo's usage.
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

	factory, err := depot.Open("sqlite3", "./test.db", &depot.FactoryOptions{
		LogSQL: true,
	})
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

	_, err = db.Exec(`
create table messages (
	id varchar primary key, 
	text varchar not null, 
	order_index integer not null, 
	len float not null, 
	attachment blob not null, 
	created timestamp not null, 
	updated timestamp
)`)
	if err != nil {
		log.Fatal(err)
	}
}
