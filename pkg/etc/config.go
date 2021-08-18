package etc

type Config struct {
	DownloadDir     string
	Url             string
	ThreadNum       int
	Urls            []string
	filterTypes     string
	filterDateTypes []string
}
