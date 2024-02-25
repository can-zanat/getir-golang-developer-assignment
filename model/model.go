package model

import (
	"errors"
	"time"
)

type GetInfoRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

type GetInfoResponse struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Records []Record `json:"records"`
}

type Record struct {
	Key        string    `json:"key"`
	CreatedAt  time.Time `json:"createdAt"`
	TotalCount int       `json:"totalCount"`
}

type InfoData struct {
	ID        string    `json:"id" bson:"_id"`
	Key       string    `json:"key" bson:"key"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Counts    []int     `json:"counts" bson:"counts"`
}

type DbResponse struct {
	Records []Record `json:"records"`
}

func (d *GetInfoRequest) ValidateDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func (d *GetInfoRequest) Validate() error {

	if !d.ValidateDate(d.StartDate) {
		return errors.New("startDate is not in the valid format YYYY-MM-DD")
	}

	if !d.ValidateDate(d.EndDate) {
		return errors.New("endDate is not in the valid format YYYY-MM-DD")
	}

	return nil
}
