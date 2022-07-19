package file

import (
	"bytes"
	"crypto/md5"
	"os"
)

/*
2つのファイルのハッシュが同一か検証する関数
*/
func IsEqualHash(file1 *os.File, file2 *os.File) (bool, error) {
	buf1 := bytes.NewBuffer([]byte{})
	_, err := buf1.ReadFrom(file1)
	if err != nil {
		return false, nil
	}

	buf2 := bytes.NewBuffer([]byte{})
	_, err = buf2.ReadFrom(file2)
	if err != nil {
		return false, nil
	}

	hash1 := md5.Sum(buf1.Bytes())
	hash2 := md5.Sum(buf2.Bytes())

	return string(hash1[:]) == string(hash2[:]), nil
}
