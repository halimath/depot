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

	_ "github.com/mattn/go-sqlite3"
)

func TestPackage(t *testing.T) {
	prepareTestDB(t)

	cols := Cols("id", "text", "attachment")

	factory, err := Open("sqlite3", "./test-package.db")
	if err != nil {
		t.Fatal(err)
	}
	defer factory.Close()

	ctx := context.Background()
	session, ctx, err := factory.Session(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := session.Commit(); err != nil {
			t.Fatal(err)
		}
	}()

	count, err := session.QueryCount(From("messages"))
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Errorf("expected 2 messages but got %d", count)
	}

	msg, err := session.QueryOne(cols, From("messages"), Where("id", "1"))
	if err != nil {
		t.Fatal(err)
	}

	if msg["id"] != "1" || msg["text"] != "hello, world" {
		t.Errorf("expected prepared message but got: %#v\n", msg)
	}

	msgs, err := session.QueryMany(cols, From("messages"), OrderBy("id", false))
	if err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 2 {
		t.Errorf("expected 2 messages but got %d", len(msgs))
	}
	if msgs[0]["id"] != "2" {
		t.Errorf("expected first id to be 2 but got: %s", msgs[0]["id"])
	}

	err = session.InsertOne(Into("messages"), Values{"id": "3", "text": "hello, once more", "attachment": []byte{'a', 'b', 'c'}})
	if err != nil {
		t.Fatal(err)
	}

	msgs, err = session.QueryMany(cols, From("messages"), OrderBy("id", false))
	if err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 3 {
		t.Errorf("expected 3 messages but got %d", len(msgs))
	}
	if msgs[0]["text"] != "hello, once more" {
		t.Errorf("got unexpected first text: %s", msgs[0]["text"])
	}

	err = session.UpdateMany(Table("messages"), Values{"text": "hello, one more time"}, Where("id", "3"))
	if err != nil {
		t.Fatal(err)
	}

	msgs, err = session.QueryMany(cols, From("messages"), OrderBy("id", false))
	if err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 3 {
		t.Errorf("expected 3 messages but got %d", len(msgs))
	}
	if msgs[0]["text"] != "hello, one more time" {
		t.Errorf("got unexpected first text: %s", msgs[0]["text"])
	}

	err = session.DeleteMany(Table("messages"), In("id", "2", "3"))
	if err != nil {
		t.Fatal(err)
	}

	// Now use a nested transaction
	func() {
		session, _, err := factory.Session(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer session.Commit()

		msgs, err = session.QueryMany(cols, From("messages"), OrderBy("id", false))
		if err != nil {
			t.Fatal(err)
		}

		if len(msgs) != 1 {
			t.Errorf("expected 1 messages but got %d", len(msgs))
		}
		if msgs[0]["id"] != "1" {
			t.Errorf("got unexpected first id: %s", msgs[0]["id"])
		}
	}()
}

func prepareTestDB(t *testing.T) {
	os.Remove("./test-package.db")

	db, err := sql.Open("sqlite3", "./test-package.db")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("create table messages (id varchar primary key, text varchar, attachment blob)")
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
