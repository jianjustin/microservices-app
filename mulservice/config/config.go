package config

var (
	Addrs       = []string{":50055", ":50056"}
	ServiceName = "mulservice"
	HTTP_ADDR   = map[string]string{
		":50055": ":8055",
		":50056": ":8056",
	}
)
