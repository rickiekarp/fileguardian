## Usage examples

### Read storage data

Read file from storage and list it's target file
```
go run main.go foo.txt
```

Read file from storage and list it's source file
```
go run main.go 8c1390e7cc93b3796752d4e800b415fe.fgd
```

Read file from storage and prints the encryption command
```
go run main.go -e -r=recipient@host foo.txt
```

Read file from storage and prints the decryption command
```
go run main.go -d 8c1390e7cc93b3796752d4e800b415fe.fgd
```