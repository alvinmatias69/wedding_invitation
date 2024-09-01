package entities

type Config struct {
	Port            string
	StaticWebDir    string
	HiddenImagePath string
	FinalImagePath  string
	Username        string
	Password        string

	HiddenImageFilePath string
	FqIfdPath           string
	HiddenImageUrlTag   string
	HiddenImageUrlValue string
	HiddenImageTokenTag string

	JwtKey      string
	JwtExpiryMs uint64
}
