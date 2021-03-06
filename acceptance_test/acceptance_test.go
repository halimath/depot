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

package acceptancetest

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/halimath/depot"
	"github.com/halimath/depot/engine/mysql"
	"github.com/halimath/depot/engine/postgres"
	"github.com/halimath/depot/engine/sqlite"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

//go:generate depot generate-repo --table=messages --out ./messagerepo_gen.go $GOFILE Message

type (
	// Message demonstrates a persistent struct showing several mapped fields.
	Message struct {
		ID         string     `depot:"id,id"`
		Text       string     `depot:"text"`
		OrderIndex int        `depot:"order_index"`
		Length     float32    `depot:"len"`
		Attachment []byte     `depot:"attachment"`
		Created    time.Time  `depot:"created"`
		Updated    *time.Time `depot:"updated,nullable"`
	}
)

func TestMariaDB(t *testing.T) {
	db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/test?parseTime=true")
	defer db.Close()

	_, err := db.Exec(`drop table if exists messages`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
create table messages (
	id varchar(255) primary key, 
	text varchar(1024) not null, 
	order_index int not null, 
	len float not null, 
	attachment blob not null, 
	created timestamp not null default 0, 
	updated timestamp null
)`)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("acceptance test", func(t *testing.T) {
		runTest(t, db, depot.Options{
			Dialect: &mysql.Dialect{},
		})
	})
}

func TestPostgres(t *testing.T) {
	db, _ := sql.Open("postgres", "host=localhost port=5432 user=user password=password dbname=test sslmode=disable")
	defer db.Close()

	_, err := db.Exec(`drop table if exists messages`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
create table messages (
	id varchar(255) primary key, 
	text varchar(1024) not null, 
	order_index int not null, 
	len float not null, 
	attachment bytea not null, 
	created timestamp not null, 
	updated timestamp null
)`)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("acceptance test", func(t *testing.T) {
		runTest(t, db, depot.Options{
			Dialect: &postgres.Dialect{},
		})
	})
}

func TestSQLite(t *testing.T) {
	const dbFile = "./test.db"
	os.Remove(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	_, err = db.Exec(`
create table messages (
	id varchar(255) primary key, 
	text varchar(1024) not null, 
	order_index int not null, 
	len float not null, 
	attachment blob not null, 
	created timestamp not null default 0, 
	updated timestamp null
)`)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("acceptance test", func(t *testing.T) {
		runTest(t, db, depot.Options{
			Dialect: &sqlite.Dialect{},
		})
	})
}

func runTest(t *testing.T, pool *sql.DB, opts depot.Options) {
	db := depot.New(pool, opts)
	defer db.Close()

	repo := &MessageRepo{
		db: db,
	}

	ctx := context.Background()
	ctx, err := repo.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := repo.Commit(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	want := Message{
		ID:         "1",
		Text:       "hello, world",
		Length:     12.0,
		Attachment: []byte{1, 2, 3},
		Created:    time.Now().UTC().Round(time.Second),
	}

	err = repo.Insert(ctx, &want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.LoadByID(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(want, *got); diff != nil {
		t.Errorf("unexpected value when loading after insert: %s", diff)
	}

	want.Text = "hello, one more time"
	var updated = time.Now().UTC().Round(time.Second)
	want.Updated = &updated

	if err := repo.Update(ctx, &want); err != nil {
		t.Fatal(err)
	}

	got, err = repo.LoadByID(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(want, *got); diff != nil {
		t.Errorf("unexpected value when loading after update: %s", diff)
	}

	err = repo.DeleteByID(ctx, "1")
	if err != nil {
		t.Error(err)
	}
}
