package domain

import (
	"time"
)

const (
	ProvinceNode = "Province"
	CityNode     = "City"
	RegencyNode  = "Regency"
	DistrictNode = "District"
	VillageNode  = "Village"
)

var (
	Incoming = map[string]string{
		"country":  CountryNode,
		"province": ProvinceNode,
		"city":     CityNode,
		"regency":  RegencyNode,
		"district": DistrictNode,
		"village":  VillageNode,
		"currency": CountryNode,
	}

	Outgoing = map[string]string{
		"currencies": CurrencyNode,
	}
)

type (
	Region struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
		Regions
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	Regions struct {
		Provinces []*Region `json:"provinces,omitempty"`
		Cities    []*Region `json:"cities,omitempty"`
		Regencies []*Region `json:"regencies,omitempty"`
		Districts []*Region `json:"districts,omitempty"`
		Villages  []*Region `json:"villages,omitempty"`
	}
)
