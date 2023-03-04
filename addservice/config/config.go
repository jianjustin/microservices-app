package config

var (
	Addrs       = []string{":50051", ":50052"}
	ServiceName = "addservice"
	HTTP_ADDR   = map[string]string{
		":50051": ":8051",
		":50052": ":8052",
	}
)
