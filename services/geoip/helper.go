package geoip

type Helper struct{}

func NewHelper() *Helper {
	return &Helper{}
}

func (h *Helper) IsAsia(data *Data) bool {
	if data == nil {
		return false
	}
	return data.Continent == "AS"
}
