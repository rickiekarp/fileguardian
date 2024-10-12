package fileprocessor

type ProcessType int

const (
	Print ProcessType = iota
	Encrypt
	Decrypt
)
