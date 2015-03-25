package metrics

import (
	"log"
	"time"
)

// Output each metric in the given registry periodically using the given
// logger.
func Log(r Registry, d time.Duration, l *log.Logger) {
	for _ = range time.Tick(d) {
		r.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case Counter:
				l.Printf("counter %s\n", name)
				l.Printf("  count:       %9d\n", metric.Count())
			case Gauge:
				l.Printf("gauge %s\n", name)
				l.Printf("  value:       %9d\n", metric.Value())
			case GaugeFloat64:
				l.Printf("gauge %s\n", name)
				l.Printf("  value:       %f\n", metric.Value())
			case Healthcheck:
				metric.Check()
				l.Printf("healthcheck %s\n", name)
				l.Printf("  error:       %v\n", metric.Error())
			case Histogram:
				h := metric.Snapshot()
				ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				l.Printf("histogram %s\n", name)
				l.Printf("  count:       %9d\n", h.Count())
				l.Printf("  min:         %9d\n", h.Min())
				l.Printf("  max:         %9d\n", h.Max())
				l.Printf("  mean:        %12.2f\n", h.Mean())
				l.Printf("  stddev:      %12.2f\n", h.StdDev())
				l.Printf("  median:      %12.2f\n", ps[0])
				l.Printf("  75%%:         %12.2f\n", ps[1])
				l.Printf("  95%%:         %12.2f\n", ps[2])
				l.Printf("  99%%:         %12.2f\n", ps[3])
				l.Printf("  99.9%%:       %12.2f\n", ps[4])
			case Meter:
				m := metric.Snapshot()
				l.Printf("meter %s\n", name)
				l.Printf("  count:       %9d\n", m.Count())
				l.Printf("  1-min rate:  %12.2f\n", m.Rate1())
				l.Printf("  5-min rate:  %12.2f\n", m.Rate5())
				l.Printf("  15-min rate: %12.2f\n", m.Rate15())
				l.Printf("  mean rate:   %12.2f\n", m.RateMean())
			case Timer:
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				l.Printf("timer %s\n", name)
				l.Printf("  count:       %9d\n", t.Count())
				l.Printf("  min:         %9d\n", t.Min())
				l.Printf("  max:         %9d\n", t.Max())
				l.Printf("  mean:        %12.2f\n", t.Mean())
				l.Printf("  stddev:      %12.2f\n", t.StdDev())
				l.Printf("  median:      %12.2f\n", ps[0])
				l.Printf("  75%%:         %12.2f\n", ps[1])
				l.Printf("  95%%:         %12.2f\n", ps[2])
				l.Printf("  99%%:         %12.2f\n", ps[3])
				l.Printf("  99.9%%:       %12.2f\n", ps[4])
				l.Printf("  1-min rate:  %12.2f\n", t.Rate1())
				l.Printf("  5-min rate:  %12.2f\n", t.Rate5())
				l.Printf("  15-min rate: %12.2f\n", t.Rate15())
				l.Printf("  mean rate:   %12.2f\n", t.RateMean())
			}
		})
	}
}

// Output each metric in the given registry periodically using the given
// logger, using a more compact display of 1 line per metric, and all times in ms.
func LogCompact(r Registry, d time.Duration, l *log.Logger) {
	for _ = range time.Tick(d) {
		r.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case Counter:
				l.Printf("count %20s,count:%9d\n", name, metric.Count())
			case Gauge:
				l.Printf("gauge %20s,value:%9d\n", name, metric.Value())
			case GaugeFloat64:
				l.Printf("gauge %20s,value:%9f\n", name, metric.Value())
			case Healthcheck:
				metric.Check()
				l.Printf("check %20s,error: %v\n", name, metric.Error())
			case Histogram:
				h := metric.Snapshot()
				ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				l.Printf("hgram %20s,count:%9d,min:%4.2f,mean:%4.2f,95%%:%4.2f,99%%:%5.2f,max:%5.2f\n",
					name, h.Count(),
					time.Duration(h.Min()).Seconds(),
					time.Duration(h.Mean()).Seconds(),
					time.Duration(ps[2]).Seconds(),
					time.Duration(ps[3]).Seconds(),
					time.Duration(h.Max()).Seconds(),
				)
			case Meter:
				m := metric.Snapshot()
				l.Printf("meter %20s,count:%9d,1mrate:%12.2f,5mrate:%12.2f,meanrate:%12.2f/s\n",
					name, m.Count(), m.Rate1(), m.Rate5(), m.RateMean())
			case Timer:
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				l.Printf("timer %20s,count:%9d,meanrate:%4.2f/s,min:%4.2f,mean:%4.2f,95%%:%4.2f,99%%:%5.2f,max:%5.2f\n",
					name, t.Count(),
					t.RateMean(),
					time.Duration(t.Min()).Seconds(),
					time.Duration(t.Mean()).Seconds(),
					time.Duration(ps[2]).Seconds(),
					time.Duration(ps[3]).Seconds(),
					time.Duration(t.Max()).Seconds(),
				)
			}
		})
	}
}
