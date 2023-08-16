package cron

import (
	"context"
	"github.com/robfig/cron/v3"
)

type Cron struct {
	*cron.Cron
}

func (c *Cron) Provide(context.Context) interface{} {
	croner := &Cron{
		Cron: cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow |
				cron.Descriptor,
		))),
	}
	croner.Start()
	return croner
}
