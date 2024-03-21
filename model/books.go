package model

type Book struct {
	ID       uint64 `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string `bson:"name" json:"name"`
	AuthorID uint64 `json:"author_id" bson:"author_id"`
}
