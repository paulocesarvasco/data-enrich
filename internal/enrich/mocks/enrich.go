package mocks

type enrich struct {
	countrySearchError error
	regionError        error
	regionResponse     string
	countryResponse    string
}

func NewEnrichService() *enrich {
	return &enrich{}
}

func (e *enrich) SetCountryResponse(country string) {
	e.countryResponse = country
}

func (e *enrich) SetCountrySearchError(err error) {
	e.countrySearchError = err
}

func (e *enrich) SetRegionResponse(resgion string) {
	e.regionResponse = resgion
}

func (e *enrich) SetSearchRegionError(err error) {
	e.regionError = err
}

func (e *enrich) GetCountryFromIp(ip string) (string, error) {
	return e.countryResponse, e.countrySearchError
}

func (e *enrich) GetCountryRegion(country string) (string, error) {
	return e.regionResponse, e.regionError
}
