package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// Playlist holds the schema definition for the Playlist entity.
type Playlist struct {
	ent.Schema
}

// Fields of the Playlist.
func (Playlist) Fields() []ent.Field {
	return []ent.Field{
		field.String("ID").Unique().Immutable().DefaultFunc(func() string {
			id, _ := gonanoid.New()
			return id
		}),
		field.String("name"),
		field.String("description"),
		field.String("imageURL"),
		field.String("playlistID"),
		field.Bool("shouldBeBackedUp"),
	}
}

// Edges of the Playlist.
func (Playlist) Edges() []ent.Edge {
	return nil
}
