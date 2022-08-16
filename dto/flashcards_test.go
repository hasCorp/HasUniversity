package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMapToCardSet(t *testing.T) {
	m := defaultMap(t)
	card := mapToCardSet(m)
	assert.Equal(t, m["uuid"], card.UUID)
	assert.Equal(t, m["author"], card.Author)
}

func TestMapToCardSetMissing(t *testing.T) {
	attributes := []string{}
	for k := range defaultMap(t) {
		attributes = append(attributes, k)
	}

	for _, attr := range attributes {
		t.Run(attr, func(t *testing.T) {
			m := defaultMap(t)
			v, ok := m[attr]
			assert.True(t, ok)
			assert.NotNil(t, mapToCardSet(m))
			delete(m, attr)
			assert.Nil(t, mapToCardSet(m))
			m[attr] = v
			assert.NotNil(t, mapToCardSet(m))
		})
	}
}

func defaultMap(t *testing.T) map[string]interface{} {
	t.Helper()
	return map[string]interface{}{
		"uuid":   "abc",
		"author": "hasanabi",
		"name":   "I'm saying it",
		"tags":   []string{"Twitch"},
		"cards": []map[string]interface{}{
			{
				"front": "peepoHas",
				"back":  "blammo",
			},
		},
		"last_update_timestamp": time.Date(2022, time.August, 16, 0, 0, 0, 0, time.Local),
	}
}
