## Usage examples

### Create a new storage

Creates a new storage file (with name)
```
go run main.go -create (-name=foo)
```

Populate storage with data
```
go run main.go -add -data=/home/user/some/path
```

### Read storage data

Read storage file and print all entries
```
go run main.go
```

Read all entries from the storage and print steps on how to decode them
```
go run main.go -instructionType=decrypt
```

Read all entries from the storage and print steps on how to encode them
```
go run main.go -instructionType=encrypt -encryptionRecipient=SOME_RECIPIENT
```

### Look up a file in the storage

Look up an encoded file
```
go run main.go -lookupDstFile=ae7b5596808cf5ccccee130a647bc19f.fgd
```

Find an encoded file by the unencoded file
```
go run main.go -lookupSrcFile=foo.bar
```