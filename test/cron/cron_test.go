package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

var a = 0
var b = 0

func TestRobfigCron(t *testing.T) {
	parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	sc, err := parser.Parse("@every 3s")
	if err != nil {
		assert.Fail(t, fmt.Sprint(err))
		return
	}

	now := time.Now()
	ret := sc.Next(now)
	fmt.Println("next=", ret)
	fmt.Println("unix=", ret.Unix())

}

func A() {
	a++
	fmt.Println("a")
}

func B() {
	b++
	fmt.Println("b")
}
