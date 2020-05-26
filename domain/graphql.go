package domain

import (
	"fmt"
	"net/http"

	"github.com/dynastymasra/cartographer/config"

	"github.com/labstack/gommon/random"

	scalar "github.com/dynastymasra/cookbook/graphql"
	"github.com/graphql-go/graphql"
)

var (
	countryField = graphql.Fields{
		"id": &graphql.Field{
			Type: scalar.UUID,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"dialCode": &graphql.Field{
			Type: graphql.String,
		},
		"ISO3166Alpha2": &graphql.Field{
			Type: graphql.String,
		},
		"ISO3166Alpha3": &graphql.Field{
			Type: graphql.String,
		},
		"ISO3166Numeric": &graphql.Field{
			Type: graphql.String,
		},
		"flags": &graphql.Field{
			Type: flagType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var country Country

				switch c := p.Source.(type) {
				case *Country:
					country = *c
				case Country:
					country = c
				}

				if err := country.Unmarshal(); err != nil {
					return nil, config.NewError(http.StatusInternalServerError, "", err.Error())
				}

				return country.Flag, nil
			},
		},
		"currencies": &graphql.Field{
			Type: graphql.NewList(CurrencyType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var currencies []*Currency

				switch country := p.Source.(type) {
				case *Country:
					currencies = country.Currencies
				case Country:
					currencies = country.Currencies
				}

				return currencies, nil
			},
		},
		"createdAt": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.DateTime,
		},
		"provinces": &graphql.Field{
			Type: graphql.NewList(RegionType(ProvinceNode)),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var provinces []*Region

				switch country := p.Source.(type) {
				case Country:
					provinces = country.Provinces
				case *Country:
					provinces = country.Provinces
				}

				return provinces, nil
			},
		},
	}

	flagType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Flag",
		Description: "Country flags",
		Fields: graphql.Fields{
			"flat": &graphql.Field{
				Type: sizeType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var size Size

					switch flag := p.Source.(type) {
					case *Flag:
						size = flag.Flat
					case Flag:
						size = flag.Flat
					}

					return size, nil
				},
			},
			"shiny": &graphql.Field{
				Type: sizeType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var size Size

					switch flag := p.Source.(type) {
					case *Flag:
						size = flag.Shiny
					case Flag:
						size = flag.Shiny
					}

					return size, nil
				},
			},
		},
	})

	sizeType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Size",
			Description: "Size of countries flags",
			Fields: graphql.Fields{
				"sixteen": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var sixteen string

						switch size := p.Source.(type) {
						case *Size:
							sixteen = size.Sixteen
						case Size:
							sixteen = size.Sixteen
						}

						return sixteen, nil
					},
				},
				"twentyFour": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var twentyFour string

						switch size := p.Source.(type) {
						case *Size:
							twentyFour = size.TwentyFour
						case Size:
							twentyFour = size.TwentyFour
						}

						return twentyFour, nil
					},
				},
				"thirtyTwo": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var thirtyTwo string

						switch size := p.Source.(type) {
						case *Size:
							thirtyTwo = size.ThirtyTwo
						case Size:
							thirtyTwo = size.ThirtyTwo
						}

						return thirtyTwo, nil
					},
				},
				"fortyEight": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var fortyEight string

						switch size := p.Source.(type) {
						case *Size:
							fortyEight = size.FortyEight
						case Size:
							fortyEight = size.FortyEight
						}

						return fortyEight, nil
					},
				},
				"sixtyFour": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var sixtyFour string

						switch size := p.Source.(type) {
						case *Size:
							sixtyFour = size.SixtyFour
						case Size:
							sixtyFour = size.SixtyFour
						}

						return sixtyFour, nil
					},
				},
			},
		})

	currencyField = graphql.Fields{
		"id": &graphql.Field{
			Type: scalar.UUID,
		},
		"ISO4217Name": &graphql.Field{
			Type: graphql.String,
		},
		"ISO4217Alphabetic": &graphql.Field{
			Type: graphql.String,
		},
		"ISO4217Numeric": &graphql.Field{
			Type: graphql.String,
		},
		"ISO4217MinorUnit": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.DateTime,
		},
	}

	RegionArgs = graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: scalar.UUID,
		},
		"code": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	}

	ListRegionArgs = graphql.FieldConfigArgument{
		"code": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"limit": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: config.Limit,
		},
		"offset": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: config.Offset,
		},
		"province": &graphql.ArgumentConfig{
			Type: RegionInput(ProvinceNode),
		},
		"city": &graphql.ArgumentConfig{
			Type: RegionInput(CityNode),
		},
		"regency": &graphql.ArgumentConfig{
			Type: RegionInput(RegencyNode),
		},
		"district": &graphql.ArgumentConfig{
			Type: RegionInput(DistrictNode),
		},
		"country": &graphql.ArgumentConfig{
			Type: CountryInput,
		},
	}

	regionFields = graphql.Fields{
		"id": &graphql.Field{
			Type: scalar.UUID,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"code": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.DateTime,
		},
		"provinces": &graphql.Field{
			Type: graphql.NewList(RegionType(ProvinceNode)),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var provinces []*Region

				switch region := p.Source.(type) {
				case *Region:
					provinces = region.Provinces
				case Region:
					provinces = region.Provinces
				}

				return provinces, nil
			},
		},
		"cities": &graphql.Field{
			Type: graphql.NewList(RegionType(CityNode)),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var cities []*Region

				switch region := p.Source.(type) {
				case *Region:
					cities = region.Cities
				case Region:
					cities = region.Cities
				}

				return cities, nil
			},
		},
		"regencies": &graphql.Field{
			Type: graphql.NewList(RegionType(RegencyNode)),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var regencies []*Region

				switch region := p.Source.(type) {
				case *Region:
					regencies = region.Regencies
				case Region:
					regencies = region.Regencies
				}

				return regencies, nil
			},
		},
		"districts": &graphql.Field{
			Type: graphql.NewList(RegionType(DistrictNode)),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var districts []*Region

				switch region := p.Source.(type) {
				case *Region:
					districts = region.Districts
				case Region:
					districts = region.Districts
				}

				return districts, nil
			},
		},
		"villages": &graphql.Field{
			Type: graphql.NewList(RegionType(VillageNode)),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var villages []*Region

				switch region := p.Source.(type) {
				case *Region:
					villages = region.Villages
				case Region:
					villages = region.Villages
				}

				return villages, nil
			},
		},
	}

	ProvinceType = graphql.NewObject(graphql.ObjectConfig{
		Name:        ProvinceNode,
		Description: "Province administrative division",
		Fields:      regionFields,
	})

	CityType = graphql.NewObject(graphql.ObjectConfig{
		Name:        CityNode,
		Description: "City administrative division",
		Fields:      regionFields,
	})

	RegencyType = graphql.NewObject(graphql.ObjectConfig{
		Name:        RegencyNode,
		Description: "Regency administrative division",
		Fields:      regionFields,
	})

	DistrictType = graphql.NewObject(graphql.ObjectConfig{
		Name:        DistrictNode,
		Description: "District administrative division",
		Fields:      regionFields,
	})

	VillageType = graphql.NewObject(graphql.ObjectConfig{
		Name:        VillageNode,
		Description: "Village administrative division",
		Fields:      regionFields,
	})

	CurrencyType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Currency",
		Description: "Currency information with ISO 4217",
		Fields:      currencyField,
	})

	CurrencyInput = graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        fmt.Sprintf("Currency_%s", random.String(5, random.Alphabetic)),
		Description: "Currency input arguments",
		Fields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type: scalar.UUID,
			},
			"ISO4217Name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"ISO4217Alphabetic": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"ISO4217Numeric": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"ISO4217MinorUnit": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})

	CountryType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Country",
		Description: "Country information with ISO 3166",
		Fields:      countryField,
	})

	CountryInput = graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "Country",
		Description: "Country input arguments",
		Fields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type: scalar.UUID,
			},
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"dialCode": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"ISO3166Alpha2": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"ISO3166Alpha3": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"ISO3166Numeric": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"currency": &graphql.InputObjectFieldConfig{
				Type: CurrencyInput,
			},
		},
	})

	CountryArgs = graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: scalar.UUID,
		},
		"ISO3166Alpha2": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"ISO3166Alpha3": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"ISO3166Numeric": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	}

	ListCountryArgs = graphql.FieldConfigArgument{
		"limit": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: config.Limit,
		},
		"offset": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: config.Offset,
		},
		"dialCode": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"currencies": &graphql.ArgumentConfig{
			Type: graphql.NewList(CurrencyInput),
		},
	}
)

func RegionInput(name string) *graphql.InputObject {
	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        fmt.Sprintf("%s_%s", name, random.New().String(5, random.Alphabetic)),
		Description: fmt.Sprintf("Region type of %s", name),
		Fields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type: scalar.UUID,
			},
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"code": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})
}

func RegionType(name string) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        fmt.Sprintf("%s_%s", name, random.New().String(5, random.Alphabetic)),
		Description: fmt.Sprintf("Region type of %s", name),
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: scalar.UUID,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"code": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	})
}
