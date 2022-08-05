package chart

//go:generate go run main.go

import (
	"Sneakers/scraper"
	"net/http"
	"time"
	

	charts "github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func findMin(x []float64) float64 {
	min := x[0]

	for _, r := range x {
		if ( r < min ) {
			min = r
		}
	}
	
	return min
}

func findMax(x []float64) float64 {
	max := x[0]

	for _, r := range x {
		if ( r > max ) {
			max = r
		}
	}
	return max
}

func xvalues(rawx []string) []time.Time {
	var dates []time.Time
	for _, ts := range rawx {
		parsed, _ := time.Parse(charts.DefaultDateFormat, ts)
		dates = append(dates, parsed)
	}
	return dates
}



func ChartIt(res http.ResponseWriter, req *http.Request, name string, size string) {
	yv, xv := scraper.GetPrices(name, size)
	dateSold := xvalues(xv)
	
	priceSeries := charts.TimeSeries{
		Name: "Price History",
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
