package geoip

type Helper struct{}

func NewHelper() *Helper {
	return &Helper{}
}

func (h *Helper) InCountries(data *Data, countries ...string) bool {
	if data == nil {
		return false
	}
	for _, country := range countries {
		if data.Country == country {
			return true

		}
	}
	return false
}
