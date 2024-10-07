package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/config/local"
)

type Test struct {
	Hello string
}

func TestLocalConfig(t *testing.T) {
	lProvider := local.New(
		local.WithConfigName("config.yaml"),
		local.WithConfigPath("./"),
		local.WithConfigType("yaml"),
	)
	config.New()
	config.SetProvider(lProvider)
	test := &Test{}
	if err := lProvider.ReadConfig(); err != nil {
		t.Error(err.Error())
	}
	if err := lProvider.Unmarshal("test", test); err != nil {
		t.Error(err.Error())
	}
	if err := lProvider.OnWatch(func() {
		if err := lProvider.Unmarshal("test", test); err != nil {
			t.Error(err.Error())
		}
	}); err != nil {
		t.Error(err.Error())
	}
	for {
		time.Sleep(time.Second)
		fmt.Println(test.Hello)
	}
}
