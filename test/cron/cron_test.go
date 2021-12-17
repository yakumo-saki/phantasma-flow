package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/huandu/go-assert"
	"github.com/robfig/cron/v3"
)

var a = 0
var b = 0

func TestRobfigCron(t *testing.T) {
	a := assert.New(t)

	parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	sc, err := parser.Parse("@every 3s")
	a.NonNilError(err)

	now := time.Now()
	ret := sc.Next(now)
	fmt.Println("next=", ret)
	fmt.Println("unix=", ret.Unix())

}
