package main

import (
	"flag"
	"github.com/efimovad/avito-internship2/internal/app"
	"time"
)

func main() {
	var cidr = flag.Int("cidr", 24, "enter mask")
	var limit = flag.Int("limit", 100, "enter requests limit per period")
	var period = flag.Duration("period", time.Minute, "enter period in seconds")
	//period := time.Second * (*periodInSec)
	var wait = flag.Duration("wait", 2 * time.Minute, "enter wait in seconds")
	//wait := time.Second * (*waitInSec)
    var port = flag.Int("port", 8080, "enter port")

	flag.Parse()
	s := app.NewServer(*cidr, *limit, *period, *wait)
	_ = s.Start(*port)
}