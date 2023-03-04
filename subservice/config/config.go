package config

var (
	Addrs       = []string{":50053", ":50054"}
	ServiceName = "subservice"
	HTTP_ADDR   = map[string]string{
		":50053": ":8053",
		":50054": ":8054",
	}
)
