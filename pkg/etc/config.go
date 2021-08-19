package etc

type Config struct {
	DownloadDir     string
	Url             string
	ThreadNum       int
	Urls            []string
	FilterTypes     string
	FilterDateTypes []string

}
