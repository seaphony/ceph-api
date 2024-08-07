package api

type Config struct {
	HttpPort       int  `yaml:"httpPort"`
	GrpcPort       int  `yaml:"grpcPort"`
	GrpcReflection bool `yaml:"grpcReflection"`
	Secure         bool `yaml:"secure"`
	ServeDebug     bool `yaml:"serveDebug"`
	AccessLog      bool `yaml:"accessLog"`
}
