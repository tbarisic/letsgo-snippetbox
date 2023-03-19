package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB  *sql.DB
	LOG *log.Logger
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `
				INSERT INTO snippets(title, content, created, expires) 
				VALUES($1, $2, now(), now() + $3::INTERVAL) RETURNING id
			`

	lastInsertId := 0
	err := m.DB.QueryRow(stmt, title, content, fmt.Sprintf("%d day", expires)).Scan(&lastInsertId)

	if err != nil {
		return 0, err
	}

	return lastInsertId, nil

}

func (m *SnippetModel) Get(id int) (*Snippet, error) {

	stmt := `
				SELECT id, title, content, created, expires
				FROM snippets
				WHERE expires >  now() AND id = $1
			`

	s := &Snippet{}

	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {

	stmt := `
				SELECT id, title, content, created, expires
				FROM snippets
				WHERE expires > now()
				ORDER BY id DESC
				LIMIT 10
			`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
