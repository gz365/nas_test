package conf

type Conf struct {
	Host     string
	Username string
	Password string
}

func GetConf() Conf {
	return Conf{
		Host:     "10.2.230.25:445",
		Username: "13521898060",
		Password: "1234",
	}
}
