package model

type User struct {
	ID            uint64 `bson:"_id,omitempty" json:"id,omitempty"`
	Login         string `bson:"login" json:"login"`
	Password      string `bson:"password" json:"password"`
	Role          string `bson:"role" json:"role"`
	RecoveryToken string `bson:"recoveryToken" json:"recoveryToken"`
}
