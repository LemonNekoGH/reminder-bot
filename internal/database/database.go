package database

import (
	"context"
	"log"

	"entgo.io/ent/dialect"
	"github.com/lemonnekogh/reminderbot/config"
	"github.com/lemonnekogh/reminderbot/ent"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

type DatabaseServiceParam struct {
	fx.In

	Config *config.Config
}

type DatabaseService struct {
	Client *ent.Client
}

func NewDatabaseService() func(DatabaseServiceParam) (*DatabaseService, error) {
	return func(dsp DatabaseServiceParam) (*DatabaseService, error) {
		client, err := ent.Open(dialect.Postgres, dsp.Config.DatabaseConnectUri)
		if err != nil {
			return nil, err
		}

		if err := client.Schema.Create(context.Background()); err != nil {
			log.Printf("数据库迁移失败： %v\n", err)
		}

		return &DatabaseService{
			Client: client,
		}, err
	}
}
