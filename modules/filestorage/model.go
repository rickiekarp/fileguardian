package filestorage

type File struct {
	ProcessId string
	Type      string
	Src       string
	Dst       string
	CreatedAt int64
}
