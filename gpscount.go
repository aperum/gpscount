package gpscount

import (
	"errors"
	"sync"

	"github.com/stratoberry/go-gpsd"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type GPSCount struct {
	Url	string

	done	chan bool
	gps	*gpsd.Session
	wg	sync.WaitGroup
}

var sampleConfig = `
  # gpsd daemon listening adress to connect to
  url = "localhost:2947"
`

func (g *GPSCount) Description() string {
	return "Read satellite numbers seen and used from gpsd daemon"
}

func (g *GPSCount) SampleConfig() string {
	return sampleConfig
}

func (g *GPSCount) Gather(acc telegraf.Accumulator) error {
	return nil
}

func (g *GPSCount) Start(acc telegraf.Accumulator) error {
	var err error

	g.gps, err = g.createGPSDSession(g.Url)
	if err != nil {
		return err
	}

	skyfilter := func(r interface{}) {
		var used, visible int

	        sky := r.(*gpsd.SKYReport)
	        for _, sat := range sky.Satellites {
	                if sat.Used {
	                        used++
	                }
	        }
	        visible = len(sky.Satellites)

		fields := map[string]interface{}{
			"visible":	visible,
			"used":		used,
		}
		acc.AddCounter("gpscount", fields, nil)
	}

	g.gps.AddFilter("SKY", skyfilter)
		
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
        	g.done = g.gps.Watch()
		<- g.done
	}()

	return nil
}

func (g *GPSCount) Stop() {
	g.done <- true
	g.wg.Wait()
}

func init() {
	inputs.Add("gpscount", func() telegraf.Input { return &GPSCount{} })
}

func (g *GPSCount) createGPSDSession(url string) (*gpsd.Session, error) {
	return gpsd.Dial(url)
}
