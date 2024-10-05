package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable().DefaultFunc(func() string {
			id, _ := gonanoid.New()
			return id
		}),
		field.String("username"),
		field.String("password"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
