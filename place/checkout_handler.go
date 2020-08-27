package place

import (
	"encoding/json"
	"net/http"

	"gitdev.inno.ktb/coach/thaichanabe/log"
	"github.com/jinzhu/gorm"
)

type checkouter interface {
	CheckOut(string) error
}

type CheckinID struct {
	ID uint
}

func CheckOutHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.Parse(r.Context())
		visiter := NewUnPersistCheckIn(db)

		db.SetLogger(log.GormLogger{Logger: logger})

		var checkIn CheckinID
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

		if err := visiter.CheckOut(checkIn.ID); err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&Message{
				Code:     "2001",
				Response: err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(&Message{
			Code: "0000",
		})
	}
}
