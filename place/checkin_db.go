package place

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type CheckInFunc func(CheckIn) (uint, error)

func (fn CheckInFunc) CheckIn(c CheckIn) (uint, error) {
	return fn(c)
}

// NewPersistCheckIn returns persist func to check-in info to db
func NewPersistCheckIn(db *gorm.DB) CheckInFunc {
	return func(chkIn CheckIn) (uint, error) {
		if err := db.Create(&chkIn).Error; err != nil {
			return 0, errors.Wrap(err, "insert check-in")
		}

		return chkIn.ID, nil
	}
}

type CheckOutFunc func(uint) error

func (fn CheckOutFunc) CheckOut(id uint) error {
	return fn(id)
}

// NewUnPersistCheckIn returns unpersist func to delete check-in from db
func NewUnPersistCheckIn(db *gorm.DB) CheckOutFunc {
	return func(id uint) error {
		return errors.Wrap(
			db.Unscoped().Delete(&CheckIn{
				ID: uint(id),
			}).Error,
			"delete check-in")
	}
}

type CheckInListFunc func(string) ([]CheckIn, error)

func (fn CheckInListFunc) List(mobileNo string) ([]CheckIn, error) {
	return fn(mobileNo)
}

// NewQueryCheckIn returns unpersist func to delete check-in from db
func NewQueryCheckIn(db *gorm.DB) CheckInListFunc {
	return func(mobileNo string) ([]CheckIn, error) {
		var list []CheckIn
		err := errors.Wrap(db.Where("mobile_no = ?", mobileNo).Find(&list).Error, fmt.Sprintf("query by %q", mobileNo))
		return list, err
	}
}
