package chart

//go:generate go run main.go

import (
	"Sneakers/scraper"
	//"os"
	//"fmt"
	"net/http"
	"time"
	

	charts "github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func xvalues(rawx []string) []time.Time {
	var dates []time.Time
	for _, ts := range rawx {
		parsed, _ := time.Parse(charts.DefaultDateFormat, ts)
		dates = append(dates, parsed)
	}
	return dates
}

func ChartIt(res http.ResponseWriter, req *http.Request) {
	yv, xv := scraper.GetPrices()
	dateSold := xvalues(xv)
	priceSeries := charts.TimeSeries{
		Name: "SPY",
		Style: charts.Style{
			StrokeColor: charts.GetDefaultColor(0),
		},
		XValues: dateSold,
		YValues: yv,
	}

	smaSeries := charts.SMASeries{
		Name: "SPY - SMA",
		Style: charts.Style{
			StrokeColor:     drawing.ColorRed,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: priceSeries,
	}

	bbSeries := &charts.BollingerBandsSeries{
		Name: "SPY - Bol. Bands",
		Style: charts.Style{
			StrokeColor: drawing.ColorFromHex("efefef"),
			FillColor:   drawing.ColorFromHex("efefef").WithAlpha(64),
		},
		InnerSeries: priceSeries,
	}

	graph := charts.Chart{
		XAxis: charts.XAxis{
			TickPosition: charts.TickPositionBetweenTicks,
		},
		YAxis: charts.YAxis{
			Range: &charts.ContinuousRange{
				Max: 1000.0,
				Min: 20.0,
			},
		},
		Series: []charts.Series{
			bbSeries,
			priceSeries,
			smaSeries,
		},
	}

	res.Header().Set("Content-Type", "image/png")
	graph.Render(charts.PNG, res)

}
