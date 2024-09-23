#!/bin/bash

inFile=$1
outFile=$2
gpgRecipient=$3

echo "Encrypting file: $inFile"

echo "Calculating md5sum"
md5=$(md5sum $inFile | awk '{ print $1 }')

echo "Encrypting file"
gpg --output $outFile --encrypt --recipient $gpgRecipient $inFile

echo "Decrypting file for verification"
gpg --output $inFile.verify --decrypt $outFile

md5After=$(md5sum $inFile.verify | awk '{ print $1 }')

rm $inFile.verify
echo "Verifying checksums"
if [[ $md5 != $md5After ]]; then
	echo "Checksum differs! Exiting..."
	exit 1
fi
