package faq

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Controller struct {
	db *sql.DB
}

type FaQ struct {
	Guild, Question, Answer string
}

func getController(location string) Controller {
	if _, err := os.Stat(location); err != nil {
		_, err := os.Create(location)

		if err != nil {
			panic(err)
		}
	}
	db, err := sql.Open("sqlite3", location)

	if err != nil {
		panic(err)
	} else {
		controller := Controller{db: db}
		if err = controller.init(); err != nil {
			panic(err)
		}
		return controller
	}
}

func (c Controller) init() error {
	_, err := c.db.Exec(
		"CREATE TABLE IF NOT EXISTS faq (guild_id TEXT NOT NULL, question TEXT NOT NULL, answer TEXT NOT NULL)",
	)

	return err
}

// Add a new question / answer
func (c Controller) Add(guild string, question string, answer string) error {
	statement, err := c.db.Prepare(
		"INSERT INTO faq (guild_id, question, answer) VALUES (?,?,?)",
	)

	if err != nil {
		return err
	}

	_, err = statement.Exec(guild, question, answer)

	return err
}

func (c Controller) Set(guild string, question string, answer string) error {
	_, err := c.db.Exec(
		"UPDATE faq SET answer=? WHERE question=? AND guild_id=?", answer, question, guild,
	)

	return err
}

func (c Controller) Remove(guild string, question string) error {
	_, err := c.db.Exec(
		"DELETE FROM faq WHERE guild_id=? AND question=?", guild, question,
	)

	return err
}

func (c Controller) GetAll(guild string) ([]FaQ, error) {
	var result []FaQ

	req, err := c.db.Query("SELECT * FROM faq WHERE guild_id=?", guild)

	if err != nil {
		return result, err
	}

	for req.Next() {
		faq := FaQ{
			Guild:    "",
			Question: "",
			Answer:   "",
		}

		err = req.Scan(&faq.Guild, &faq.Question, &faq.Answer)

		if err != nil {
			return result, err
		} else {
			result = append(result, faq)
		}
	}

	return result, nil
}
