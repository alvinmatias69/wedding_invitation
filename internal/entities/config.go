package entities

type Config struct {
	Port            uint64 `toml:"port"`
	StaticWebDir    string `toml:"static_web_dir"`
	HiddenImagePath string `toml:"hidden_image_path"`
	SteamTokenPath  string `toml:"steam_token_path"`
	Username        string `toml:"username"`
	Password        string `toml:"password"`

	HiddenImageFilePath string `toml:"hidden_image_file_path"`
	FqIfdPath           string `toml:"fq_ifd_path"`
	HiddenImageUrlTag   string `toml:"hidden_image_url_tag"`
	HiddenImageTokenTag string `toml:"hidden_image_token_tag"`

	JwtKey          string `toml:"jwt_key"`
	JwtExpiryMinute uint64 `toml:"jwt_expiry_minute"`

	MessageLimit uint64 `toml:"message_limit"`
}
