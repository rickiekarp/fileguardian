package filestorage

type Storage struct {
	Id           string
	CreationDate int64
	LastModified int64
	Files        map[string][]File
}

type File struct {
	ProcessId string
	Type      string
	Src       string
	Dst       string
	CreatedAt int64
}
