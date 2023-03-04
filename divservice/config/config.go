package config

var (
	Addrs       = []string{":50057", ":50058"}
	ServiceName = "divservice"
	HTTP_ADDR   = map[string]string{
		":50057": ":8057",
		":50058": ":8058",
	}
)
