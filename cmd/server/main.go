package main

import (
	"flag"
	"github.com/efimovad/avito-internship2/internal/app"
	"time"
)

func main() {
	var cidr = flag.Int("cidr", 24, "enter mask")
	var limit = flag.Int("limit", 100, "enter requests limit per period")
	var periodInSec = flag.Duration("period", 60, "enter period in seconds")
	period := time.Second * (*periodInSec)
	var waitInSec = flag.Duration("wait", 120, "enter wait in seconds")
	wait := time.Second * (*waitInSec)

	flag.Parse()
	s := app.NewServer(*cidr, *limit, period, wait)
	_ = s.Start()
}