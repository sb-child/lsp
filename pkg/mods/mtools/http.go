package mtools

import (
	"github.com/gocolly/colly"
)

func CollyCollector() *colly.Collector {
	return colly.NewCollector(
		colly.UserAgent(
			"Mozilla/5.0 (X11; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
		),
		colly.AllowURLRevisit(),
		colly.Async(true),
	)
}
func CollyCollectorSlow() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent(
			"Mozilla/5.0 (X11; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
		),
		colly.AllowURLRevisit(),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 4,
	})
	return c
}
