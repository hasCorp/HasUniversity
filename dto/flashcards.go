package dto

import (
	"context"
	"log"
	"time"

	"github.com/gocql/gocql"
)

type Card struct {
	Front string `json:"front" cql:"front"`
	Back  string `json:"back" cql:"back"`
}

type CardSet struct {
	UUID        string    `json:"uuid,omitempty" cql:"uuid"`
	Author      string    `json:"author,omitempty" cql:"author"`
	Name        string    `json:"name,omitempty" cql:"name"`
	Tags        []string  `json:"tags,omitempty" cql:"tags"`
	Cards       []Card    `json:"cards,omitempty" cql:"cards"`
	LastUpdated time.Time `json:"last_update_timestamp,omitempty" cql:"last_update_timestamp"`
}

func GetCardSet(ctx context.Context, session *gocql.Session, uuid string) (*CardSet, error) {
	m := make(map[string]interface{})
	err := session.Query(`SELECT * FROM flashcard WHERE uuid = ?`, uuid).WithContext(ctx).MapScan(m)

	if err != nil {
		log.Println("oopsie whoopsie ", err)
		return nil, err
	}

	return mapToCardSet(m), nil
}

func GetAllCardSets(ctx context.Context, session *gocql.Session, limit int) ([]*CardSet, error) {
	var items []map[string]interface{}
	var err error
	sets := []*CardSet{}
	iter := session.Query(`SELECT uuid, author, name, tags, cards, last_update_timestamp FROM flashcard`).WithContext(ctx).Iter()

	items, err = iter.SliceMap()
	if err != nil {
		log.Println("oopsie getting all records", err)
		iter.Close()

		return sets, err
	}

	for _, item := range items {
		fc := mapToCardSet(item)
		if fc != nil {
			sets = append(sets, fc)
		}
	}

	return sets, iter.Close()
}

func AddCardSet(ctx context.Context, session *gocql.Session, c *CardSet) (*CardSet, error) {
	c.UUID = gocql.TimeUUID().String()
	c.LastUpdated = time.Now()
	err := session.Query(
		`INSERT INTO flashcard (uuid, author, name, tags, cards, last_update_timestamp) VALUES (?, ?, ?, ?, ?, ?)`,
		c.UUID,
		c.Author,
		c.Name,
		c.Tags,
		c.Cards,
		c.LastUpdated,
	).WithContext(ctx).Exec()
	if err != nil {
		log.Println("failed to insert data", err)
		return nil, err
	}

	return c, nil
}

// mapToCardSet takes a map presumably as output from
// a cql query, and then transforms it into a CardSet.
// If there is a failure, return nil.
func mapToCardSet(m map[string]interface{}) *CardSet {
	uuid, ok := m["uuid"].(string)
	if !ok {
		log.Println("failed to get uuid attribute")
		return nil
	}

	author, ok := m["author"].(string)
	if !ok {
		log.Println("failed to get author attribute")
		return nil
	}

	name, ok := m["name"].(string)
	if !ok {
		log.Println("failed to get name attribute")
	}

	tags, ok := m["tags"].([]string)
	if !ok {
		log.Println("failed to get tags attribute")
	}

	lastUpdated, ok := m["last_update_timestamp"].(time.Time)
	if !ok {
		log.Println("failed to get last_update_timestamp attribute")
		return nil
	}

	rawCards, ok := m["cards"].([]map[string]interface{})
	if !ok {
		log.Println("failed to get cards attribute")
		return nil
	}
	cards := make([]Card, len(rawCards))
	for i, card := range rawCards {
		front, ok := card["front"].(string)
		if !ok {
			log.Println("failed to get front attribute for card")
			return nil
		}
		back, ok := card["back"].(string)
		if !ok {
			log.Println("failed to get back attribute for card")
			return nil
		}
		c := Card{
			Front: front,
			Back:  back,
		}
		cards[i] = c
	}

	fc := CardSet{
		UUID:        uuid,
		Author:      author,
		Name:        name,
		Tags:        tags,
		Cards:       cards,
		LastUpdated: lastUpdated,
	}

	return &fc
}
