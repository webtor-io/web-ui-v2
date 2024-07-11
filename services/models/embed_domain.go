package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type EmbedDomain struct {
	tableName struct{}  `pg:"embed_domain"`
	ID        uuid.UUID `pg:"embed_domain_id,type:uuid,pk,default:uuid_generate_v4()"`
	Domain    string
	Email     string
	Ads       bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
