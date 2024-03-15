package filestorage

import "flag"

var InstructionType = flag.String("instructionType", "", "prints compression instructions")

// gpg recipient
var EncryptionRecipient = flag.String("encryptionRecipient", "", "sets the recipient of the encrypted file")

// entry lookup flags
var LookupSrcFile = flag.String("lookupSrcFile", "", "looks up an entry by a given source file")
var LookupDstFile = flag.String("lookupDstFile", "", "looks up an entry by a given encoded file")

type File struct {
	Tag string
	Src string
	Dst string
}
