package config

type System struct {
	Env          string `mapstructure:"env" json:"env" yaml:"env"`                               // 环境值
	Addr         int    `mapstructure:"addr" json:"addr" yaml:"addr"`                            // 端口值
	DbType       string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`                   // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	OssType      string `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"`                // Oss类型
	UseRedis     bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`             // 使用redis
	SingleSignOn bool   `mapstructure:"single-signon" json:"single-signon" yaml:"single-signon"` // 单点登录 use-multipoint 多点登录拦截
	// UseMultipoint bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"` // 多点登录拦截
	LimitCountIP int `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP  int `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`

	UseCertif bool   `mapstructure:"use-certif" json:"use-certif" yaml:"use-certif"` // 使用 TLS 证书
	CertFile  string `mapstructure:"cert-file" json:"cert-file" yaml:"cert-file"`
	KeyFile   string `mapstructure:"key-file" json:"key-file" yaml:"key-file"`

	// use-netip
	UseNetIP bool `mapstructure:"use-netip" yaml:"use-netip"`
}
