package place

import (
	"encoding/json"
	"net/http"

	"gitdev.inno.ktb/coach/thaichanabe/log"
	"github.com/jinzhu/gorm"
)

type CheckIn struct {
	ID       uint    `json:"-" gorm:"PRIMARY_KEY"`
	IP       string  `json:"ipAddress"`
	MobileNo string  `json:"mobileNo" gorm:"unique;not null"`
	Lat      float64 `json:"-"`
	Long     float64 `json:"-"`
}

func CheckInHandler(db *gorm.DB, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logger := log.Parse(r.Context())
		loc := NewLocationGetter(client, "http://ip-api.com/json/", logger)
		visiter := NewPersistCheckIn(db)

		db.SetLogger(log.GormLogger{Logger: logger})

		var checkIn CheckIn
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

		resp := CheckInToLocation(visiter, loc)(checkIn)

		w.WriteHeader(resp.Code)
		json.NewEncoder(w).Encode(resp.Message)
	}
}
