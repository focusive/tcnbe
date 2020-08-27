package place

import "net/http"

type getter interface {
	Get(string) (*IpGeo, error)
}

type checkiner interface {
	CheckIn(CheckIn) (uint, error)
}

func CheckInToLocation(visiter checkiner, location getter) func(CheckIn) *Response {
	return func(checkIn CheckIn) *Response {
		loc, err := location.Get(checkIn.IP)
		if err != nil {
			return &Response{
				Code: http.StatusOK,
				Message: Message{
					Code:     "1001",
					Response: err.Error(),
				},
			}
		}

		checkIn.Lat = loc.Lat
		checkIn.Long = loc.Lon

		id, err := visiter.CheckIn(checkIn)
		if err != nil {
			return &Response{
				Code: http.StatusOK,
				Message: Message{
					Code:     "2001",
					Response: err.Error(),
				},
			}
		}

		return &Response{
			Code: http.StatusOK,
			Message: Message{
				Code:     "0000",
				Response: id,
			},
		}
	}
}
