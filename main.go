package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func main() {
	s, _ := cron.ParseStandard("@weekly")
	fmt.Println(s.Next(time.Now().Add(time.Hour * 24 * 4)))
}
