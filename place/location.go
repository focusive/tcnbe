package place

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// IpGeo replesent response from https://ip-api.com/docs/api:json
type IpGeo struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}

type GetLocationFunc func(ip string) (*IpGeo, error)

func (fn GetLocationFunc) Get(ip string) (*IpGeo, error) {
	return fn(ip)
}

func NewLocationGetter(client *http.Client, url string, logger *zap.Logger) GetLocationFunc {
	return func(ip string) (*IpGeo, error) {
		url = fmt.Sprintf("%s%s", url, ip)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.Wrap(err, fmt.Sprintf("new request %s %s", http.MethodGet, url))
		}

		res, err := client.Do(req)
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.Wrap(err, fmt.Sprintf("do %s %s", http.MethodGet, url))
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			logger.Warn(err.Error())
			return nil, errors.New(res.Status)
		}

		var ipgeo IpGeo
		err = json.NewDecoder(res.Body).Decode(&ipgeo)
		return &ipgeo, err
	}
}
