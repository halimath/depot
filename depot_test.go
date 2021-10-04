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

package depot

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/halimath/depot/query"
	_ "github.com/mattn/go-sqlite3"
)

var (
	cols = query.Cols("id", "text", "attachment")
)

func TestReading(t *testing.T) {
	prepareTestDB(t)

	db, err := Open("sqlite3", "./test-package.db", Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, ctx, err := db.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	count, err := tx.QueryCount(query.From("messages"))
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Errorf("expected 2 messages but got %d", count)
	}

	msg, err := tx.QueryOne(cols, query.From("messages"), query.Where("id", "1"))
	if err != nil {
		t.Fatal(err)
	}

	if msg["id"] != "1" || msg["text"] != "hello, world" {
		t.Errorf("expected prepared message but got: %#v\n", msg)
	}

	msgs, err := tx.QueryMany(cols, query.From("messages"), query.OrderBy("id", false))
	if err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 2 {
		t.Errorf("expected 2 messages but got %d", len(msgs))
	}
	if msgs[0]["id"] != "2" {
		t.Errorf("expected first id to be 2 but got: %s", msgs[0]["id"])
	}

	// Now use a nested transaction
	func() {
		tx, _, err := db.BeginTx(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()

		msgs, err = tx.QueryMany(cols, query.From("messages"), query.OrderBy("id", false))
		if err != nil {
			t.Fatal(err)
		}

		if len(msgs) != 2 {
			t.Errorf("expected 2 messages but got %d", len(msgs))
		}
		if msgs[0]["id"] != "2" {
			t.Errorf("got unexpected first id: %s", msgs[0]["id"])
		}
		tx.Commit()
	}()
	tx.Commit()
}

func TestInsert(t *testing.T) {
	prepareTestDB(t)

	db, err := Open("sqlite3", "./test-package.db", Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, _, err := db.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	count, err := tx.QueryCount(query.From("messages"))
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Errorf("expected 2 messages but got %d", count)
	}

	err = tx.InsertOne(query.Into("messages"), Values{"id": "3", "text": "hello, once more", "attachment": []byte{'a', 'b', 'c'}})
	if err != nil {
		t.Fatal(err)
	}

	msgs, err := tx.QueryMany(cols, query.From("messages"), query.OrderBy("id", false))
	if err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 3 {
		t.Errorf("expected 3 messages but got %d", len(msgs))
	}
	if msgs[0]["text"] != "hello, once more" {
		t.Errorf("got unexpected first text: %s", msgs[0]["text"])
	}
}

func TestUpdate(t *testing.T) {
	prepareTestDB(t)

	db, err := Open("sqlite3", "./test-package.db", Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, _, err := db.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	err = tx.UpdateMany(query.Table("messages"), Values{"text": "hello, one more time"}, query.Where("id", "2"))
	if err != nil {
		t.Fatal(err)
	}

	msgs, err := tx.QueryMany(cols, query.From("messages"), query.OrderBy("id", false))
	if err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 2 {
		t.Errorf("expected 3 messages but got %d", len(msgs))
	}
	if msgs[0]["text"] != "hello, one more time" {
		t.Errorf("got unexpected first text: %s", msgs[0]["text"])
	}
}

func TestDelete(t *testing.T) {
	prepareTestDB(t)

	db, err := Open("sqlite3", "./test-package.db", Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, _, err := db.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	count, err := tx.QueryCount(query.From("messages"))
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Errorf("expected 2 messages but got %d", count)
	}

	err = tx.DeleteMany(query.From("messages"))
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	count, err = tx.QueryCount(query.From("messages"))
	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Errorf("expected 0 messages but got %d", count)
	}
}

func TestNullValues(t *testing.T) {
	prepareTestDB(t)

	db, err := Open("sqlite3", "./test-package.db", Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, _, err := db.BeginTx(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	err = tx.InsertOne(query.Into("messages"), Values{"id": "3", "text": "hello, once more", "attachment": nil})
	if err != nil {
		t.Fatal(err)
	}

	vals, err := tx.QueryOne(cols, query.From("messages"), query.Where("id", "3"))
	if err != nil {
		t.Fatal(err)
	}

	if vals["attachment"] != nil {
		t.Errorf("expected nil value for null column but got %#v", vals["attachment"])
	}
}

func prepareTestDB(t *testing.T) {
	os.Remove("./test-package.db")

	db, err := sql.Open("sqlite3", "./test-package.db")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("create table messages (id varchar primary key, text varchar not null, attachment blob default null)")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	_, err = db.Exec("insert into messages values (?, ?, ?)", "1", "hello, world", nil)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	_, err = db.Exec("insert into messages values (?, ?, ?)", "2", "hello, again", []byte{'a'})
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	db.Close()
}
