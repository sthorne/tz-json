package main

import (
	
	"time"
	"bufio"
	"strings"
	"net/http"
	"io/ioutil"
	"archive/tar"
	"compress/gzip"
	"encoding/json"

)

const (

	TimezoneReleases	= "http://www.iana.org/time-zones/repository/releases/"
	TimezoneFile		= "tzdata2014j.tar.gz"
	ZoneTabFile			= "zone.tab"

)

func main() {

	Zones := map[string]string{}

	// grab the latest tar
	r, e := http.Get(TimezoneReleases + TimezoneFile)

	if e != nil {
		panic(e)
	}

	defer r.Body.Close()

	// decompress
	g, e := gzip.NewReader(r.Body)

	if e != nil {
		panic(e)
	}

	defer g.Close()

	// read the archive
	t := tar.NewReader(g)

	for {
		h, e := t.Next()

		if e != nil {
			panic(e)
		}

		if h == nil {
			break
		}

		if h.Name != ZoneTabFile {
			continue
		}

		scanner := bufio.NewScanner(t)

		for scanner.Scan() {

			l := scanner.Text()

			// skip comments
			if strings.HasPrefix(l, "#") {
				continue
			}

			c := strings.Split(l, "\t")

			if len(c) < 3 {
				continue
			}

			Location, e := time.LoadLocation(c[2])

			if e != nil {
				continue
			}
			
			Zones[c[2]] = time.Now().In(Location).Format("-0700")
		}

		break
	}

	b, e := json.Marshal(Zones)

	if e != nil {
		panic(e)
	}

	e = ioutil.WriteFile("tz.json", b, 0644)

	if e != nil {
		panic(e)
	}
}