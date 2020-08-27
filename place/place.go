package place

import (
	"encoding/json"
	"net/http"

	"gitdev.inno.ktb/coach/thaichanabe/log"
	"github.com/jinzhu/gorm"
)

type Response struct {
	Code    int
	Message Message
}

type Message struct {
	Code     string
	Message  string
	Response interface{}
}

type lister interface {
	List(string) ([]CheckIn, error)
}

type CheckinMobile struct {
	MobileNo string `json:"mobileNo"`
}

func Handler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.Parse(r.Context())
		visiter := NewQueryCheckIn(db)

		db.SetLogger(log.GormLogger{Logger: logger})

		var checkIn CheckinMobile
		if err := json.NewDecoder(r.Body).Decode(&checkIn); err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Message{
				Code:    "3001",
				Message: err.Error(),
			})
			return
		}
		defer r.Body.Close()

		places, err := visiter.List(checkIn.MobileNo)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&Message{
				Code:     "2001",
				Response: err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(places)
	}
}
