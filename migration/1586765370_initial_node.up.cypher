CREATE CONSTRAINT currency_pkey ON (node:Currency) ASSERT (node.id) IS UNIQUE;
CREATE CONSTRAINT currency_iso4217_alphabetic_idx ON (node:Currency) ASSERT (node.ISO4217Alphabetic) IS UNIQUE;
CREATE CONSTRAINT currency_iso4217_numeric_idx ON (node:Currency) ASSERT (node.ISO4217Numeric) IS UNIQUE;

CREATE CONSTRAINT country_pkey ON (node:Country) ASSERT (node.id) IS UNIQUE;
CREATE CONSTRAINT country_iso3166_alpha2_idx ON (node:Country) ASSERT (node.ISO3166Alpha2) IS UNIQUE;
CREATE CONSTRAINT country_iso3166_alpha3_idx ON (node:Country) ASSERT (node.ISO3166Alpha3) IS UNIQUE;
CREATE CONSTRAINT country_iso3166_numeric_idx ON (node:Country) ASSERT (node.ISO3166Numeric) IS UNIQUE;
CREATE CONSTRAINT country_name_idx ON (node:Country) ASSERT (node.name) IS UNIQUE;

CREATE CONSTRAINT province_pkey ON (node:Province) ASSERT (node.id) IS UNIQUE;
CREATE CONSTRAINT province_code_idx ON (node:Province) ASSERT (node.code) IS UNIQUE;

CREATE CONSTRAINT city_pkey ON (node:City) ASSERT (node.id) IS UNIQUE;
CREATE CONSTRAINT city_code_idx ON (node:City) ASSERT (node.code) IS UNIQUE;

CREATE CONSTRAINT regency_pkey ON (node:Regency) ASSERT (node.id) IS UNIQUE;
CREATE CONSTRAINT regency_code_idx ON (node:Regency) ASSERT (node.code) IS UNIQUE;

CREATE CONSTRAINT district_pkey ON (node:District) ASSERT (node.id) IS UNIQUE;
CREATE CONSTRAINT district_code_idx ON (node:District) ASSERT (node.code) IS UNIQUE;

CREATE CONSTRAINT village_pkey ON (node:Village) ASSERT (node.id) IS UNIQUE;
CREATE CONSTRAINT village_code_idx ON (node:Village) ASSERT (node.code) IS UNIQUE;