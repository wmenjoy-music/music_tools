package model

type ForumArtistInfo struct {

}
type ForumAlbumInfo struct {
	Title     string
	Url       string
	Artist    string
	Name      string
	GenreType string
	BitRate   string
	FileType  string
	Duration  string
	Country   string
	Content   string
	Year      string
	MagnetLink string
	MagnetTitle string
	Torrent    string
}

type ArtistInfo struct {
	Name    string
	Id      string
	Url     string
	Memeber []string
	Country string
}

type MusicInfo struct {
	Name        string
	Id          string
	Url         string
	Album       string
	Artist      string
	Postion     string
	DataTitle   string
	DownloadUrl string
	Duration    string
	Stars       string
	BitRate     string
}

type AlbumInfo struct {
	Name       string
	Id         string
	FullName   string
	Image      string
	Artist     []ArtistInfo
	Genre      []GenreInfo
	MusicList  []MusicInfo
	DataType   string
	Year       string
	CreateDate string
	Category   string
	Url        string
	Rating     string
}

type GenreInfo struct {
	Name string
	Id   string
	Url  string
	Desc string
}

// CategoryTypeMap 枚举值
var CategoryTypeMap map[string]string = make(map[string]string, 15)
var TypeCategoryMap map[string]string = make(map[string]string, 15)

func NormalCategory(category string) string {
	return CategoryTypeMap[TypeCategoryMap[category]]
}

func init() {
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
	TypeCategoryMap["All"] = "0"
	TypeCategoryMap["Все"] = "0"
	TypeCategoryMap["Unsorted"] = "1"
	TypeCategoryMap["Тип не назначен"] = "1"
	TypeCategoryMap["Album"] = "2"
	TypeCategoryMap["Студийный альбом"] = "2"
	TypeCategoryMap["EP"] = "3"
	TypeCategoryMap["Single"] = "4"
	TypeCategoryMap["Сингл"] = "4"
	TypeCategoryMap["Bootleg"] = "5"
	TypeCategoryMap["Бутлег"] = "5"
	TypeCategoryMap["Live"] = "6"
	TypeCategoryMap["Compilation"] = "7"
	TypeCategoryMap["Сборник разных исполнителей"] = "7"
	TypeCategoryMap["MixType"] = "8"
	TypeCategoryMap["Микстейп"] = "8"
	TypeCategoryMap["Demo"] = "9"
	TypeCategoryMap["Демо"] = "9"
	TypeCategoryMap["DJ Mix"] = "10"
	TypeCategoryMap["DJ микс"] = "10"
	TypeCategoryMap["Group Compilations"] = "11"
	TypeCategoryMap["Сборник исполнителя"] = "11"
	TypeCategoryMap["Split"] = "12"
	TypeCategoryMap["Unoffical Compilation"] = "13"
	TypeCategoryMap["Неофициальный сборник"] = "13"
	TypeCategoryMap["OST"] = "14"
	TypeCategoryMap["Саундтрек"] = "14"

}
