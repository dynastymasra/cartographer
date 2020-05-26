package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	CurrencyNode = "Currency"
	CountryNode  = "Country"
)

type (
	Country struct {
		ID             string      `json:"id"`
		Name           string      `json:"name"`
		ISO3166Alpha2  string      `json:"ISO3166Alpha2"`
		ISO3166Alpha3  string      `json:"ISO3166Alpha3"`
		ISO3166Numeric string      `json:"ISO3166Numeric"`
		CallingCode    string      `json:"dialCode"`
		Currencies     []*Currency `json:"currencies"`
		Flag           Flag        `json:"flag"`
		FlagStr        string      `json:"flags"`
		Regions
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	Currency struct {
		ID                string    `json:"id"`
		ISO4217Name       string    `json:"ISO4217Name"`
		ISO4217Numeric    string    `json:"ISO4217Numeric"`
		ISO4217MinorUnit  string    `json:"ISO4217MinorUnit"`
		ISO4217Alphabetic string    `json:"ISO4217Alphabetic"`
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
	}
	Flag struct {
		Flat  Size `json:"flat"`
		Shiny Size `json:"shiny"`
	}

	Size struct {
		Sixteen    string `json:"sixteen"`
		TwentyFour string `json:"twentyFour"`
		ThirtyTwo  string `json:"thirtyTwo"`
		FortyEight string `json:"fortyEight"`
		SixtyFour  string `json:"sixtyFour"`
	}
)

func (c *Country) Unmarshal() error {
	var flag Flag

	res, err := json.Marshal(c.FlagStr)
	if err != nil {
		return err
	}

	str, err := strconv.Unquote(string(res))
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := json.Unmarshal([]byte(str), &flag); err != nil {
		return err
	}
	c.Flag = flag
	return nil
}
