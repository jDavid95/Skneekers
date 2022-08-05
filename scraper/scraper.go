package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"strconv"
	"time"
	

	"github.com/PuerkitoBio/goquery"
)

func convertDate(date string) string{
	const (
		layoutISO = "Jan 2, 2006"
		layoutUS  = "2006-01-02"
	)
	myDate, err := time.Parse(layoutISO, date)
	if err != nil {
		panic(err)
	}

	newDate := myDate.Format(layoutUS)

	return newDate
	
}

func errorCheck(error error) {
	if error != nil {
		fmt.Println(error)
	}
}

func htmlReq(url string) *http.Response {
	response, error := http.Get(url)
	errorCheck(error)
	if response.StatusCode > 400 {
		fmt.Println("Status Code:", response.StatusCode)
	}

	return response
}

func scrapeListings(doc *goquery.Document, prices *[]float64, dates *[]string) {
	doc.Find("ul.srp-results>li.s-item").Each(func(index int, item *goquery.Selection) {
		

		price_span := strings.TrimSpace(item.Find("span.s-item__price").Text())
		date_span := item.Find("div.s-item__title--tagblock>span.POSITIVE").Text()
		
		/*some listings will have 2 different prices when they are
		  sold at auctions with "Buy Now" options, so in this case
		  the variable "price"represents both the amount that
		  it was sold and the amount of the "Buy Now" option.
		This makes our data unreliable by having an extra price
		that does not represent the final sale.
		To work around this we trim the first dollar sign and then
		split the string after the second dollar sign.*/
		
		price,_ := strconv.ParseFloat(strings.Split(strings.Trim(price_span, "$"), "$")[0], 64)

		*prices = append(*(prices), price)
		
		date := convertDate(strings.Trim(date_span, "Sold "))

		*dates = append(*(dates), date)
		

		 		
	})
}

func parseQueryName(name string) string{
	return strings.Replace(name, " ", "+", -1)
}

func parseQuerySize(size string) string{
	if strings.Contains(size, ".5") == false {
		return size
	}
	return strings.Replace(size, ".5", "%252E5", -1)
}

func GetPrices(name string, size string) ([]float64, []string) {
	sneakerName := parseQueryName(name)
	sneakerSize := parseQuerySize(size)
	
	url := "https://www.ebay.com/sch/i.html?_fsrp=1&rt=nc&_from=R40&_nkw="+sneakerName+"&_sacat=0&LH_TitleDesc=0&LH_Sold=1&_oaa=1&_dcat=15709&US%2520Shoe%2520Size="+sneakerSize

	
	
	var prices []float64
	var dates []string

	for i:= 0; i < 4 ; i++ {
		response := htmlReq(url)
		defer response.Body.Close()
	
		doc, error := goquery.NewDocumentFromReader(response.Body)
		errorCheck(error)
		scrapeListings(doc, &prices, &dates)

		href, _ := doc.Find("nav.pagination>a.pagination__next").Attr("href")
	
		url = href
		
	}
	fmt.Println(dates)
	return prices, dates
		
}
