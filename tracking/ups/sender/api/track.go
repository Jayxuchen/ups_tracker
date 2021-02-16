package api

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
)

type gzreadCloser struct {
	*gzip.Reader
	io.Closer
}

type Response struct {
	TrackResponse struct {
		Shipment []struct {
			Package []struct {
				TrackingNumber string `json:"trackingNumber"`
				DeliveryDate   []struct {
					Type string `json:"type"`
					Date string `json:"date"`
				} `json:"deliveryDate"`
				DeliveryTime struct {
					StartTime string `json:"startTime"`
					EndTime   string `json:"endTime"`
					Type      string `json:"type"`
				} `json:"deliveryTime"`
				Activity []struct {
					Location struct {
						Address struct {
							City          string `json:"city"`
							StateProvince string `json:"stateProvince"`
							PostalCode    string `json:"postalCode"`
							Country       string `json:"country"`
						} `json:"address"`
					} `json:"location"`
					Status struct {
						Type        string `json:"type"`
						Description string `json:"description"`
						Code        string `json:"code"`
					} `json:"status"`
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"activity"`
			} `json:"package"`
		} `json:"shipment"`
	} `json:"trackResponse"`
}

func formatTime(time string) string {
	if len(time) != 6 {
		return time
	} else {
		formattedStr := fmt.Sprintf("%s:%s:%s", time[0:2], time[2:4], time[4:6])
		return formattedStr
	}

}

func TrackingInfo(accessLicenseNumber string, trackingNumber string, historyLen float64) *Response {

	url := fmt.Sprintf("https://onlinetools.ups.com/track/v1/details/%s", trackingNumber)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("AccessLicenseNumber", accessLicenseNumber)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer res.Body.Close()
	response := &Response{}
	json.NewDecoder(res.Body).Decode(response)

	for _, s := range response.TrackResponse.Shipment {
		//		fmt.Printf("Status: %+v\n", s.Package)
		for _, p := range s.Package {
			fmt.Printf("Tracking Number: %s\n", p.TrackingNumber)
			// print current location
			println("Current Location")
			for i := 0; i < int(math.Min(historyLen, float64(len(p.Activity)))); i++ {
				fmt.Printf("Activity: %d\n", i)
				fmt.Printf("\tStatus: %s\n", p.Activity[i].Status.Description)
				fmt.Printf("\tDate: %s-%s-%s\n", p.Activity[i].Date[4:6], p.Activity[i].Date[6:8], p.Activity[i].Date[0:4])
				fmt.Printf("\tTime: %s\n", formatTime(p.Activity[i].Time))
				fmt.Printf("\tCity: %s\n", p.Activity[i].Location.Address.City)
				fmt.Printf("\tState: %s\n", p.Activity[i].Location.Address.StateProvince)
				fmt.Printf("\tCountry: %s\n", p.Activity[i].Location.Address.Country)

			}

		}
	}
	//println(response.TrackResponse.Shipment[0].Package[0].TrackingNumber)
	return response

}

/*
func main() {

	if len(os.Args) > 3 {
		historyLen, _ := strconv.ParseFloat(os.Args[3], 64)
		trackingInfo(os.Args[1], os.Args[2], historyLen)
	} else {
		println("Usage: track <Access License Number> <Tracking Number>")
	}
}*/
