package model

type ArtistInfo struct {
	Name    string
	Id		string
	Memeber []string
	Country string
}

type MusicInfo struct {
	Name        string
	Id			string
	Url         string
	Album       string
	Artist      string
	Postion     string
	DataTitle   string
	Download    string
}

type AlbumInfo struct {
	Name       string
	Id		   string
	FullName   string
	Image      string
	Artist []ArtistInfo
	Genre      []GenreInfo
	MusicList  []MusicInfo
	DataType   string
	CreateDate string
	Category   string
	Url        string
}

type GenreInfo struct {
	Name   string
	Id     string
	Url    string
	Desc   string

}

// CategoryTypeMap 枚举值
var CategoryTypeMap map[string]string = make(map[string]string, 15)

func init(){
	CategoryTypeMap["0"] = "All"
	CategoryTypeMap["1"] = "Unsorted"
	CategoryTypeMap["2"] = "Album"
	CategoryTypeMap["3"] = "EP"
	CategoryTypeMap["4"] = "Single"
	CategoryTypeMap["5"] = "Bootleg"
	CategoryTypeMap["6"] = "Live"
	CategoryTypeMap["7"] = "Compilation"
	CategoryTypeMap["8"] = "MixType"
	CategoryTypeMap["9"] = "Demo"
	CategoryTypeMap["10"] = "DJ Mix"
	CategoryTypeMap["11"] = "Group Compilations"
	CategoryTypeMap["12"] = "Split"
	CategoryTypeMap["13"] = "Unoffical Compilation"
	CategoryTypeMap["14"] = "OST"
}


