package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

/*
ファイルをコピーする関数

対象がディレクトリだった場合はエラー

引数
target - コピー元のパス
destination - コピー先のパス

返り値
err
*/
func Copy(target string, destination string) error {
  // コピー元がファイルか検証
	targetFileInfo, err := os.Stat(target)
	if err != nil {
		return err
	}
	if targetFileInfo.IsDir() {
		return fmt.Errorf("error: %s is directory", target)
	}

  // コピー先の親ディレクトリまで生成
	err = os.MkdirAll(filepath.Dir(target), 0777)
	if err != nil {
		return err
	}

  // ファイルのコピー処理
	destinationFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	targetFile, err := os.Open(target)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(destinationFile, targetFile)
	if err != nil {
		return err
	}

  // ファイルのハッシュがコピー前後で同一か検証
	bool, err := IsEqualHash(destinationFile, targetFile)
	if bool {
		return nil
	}

	err_ := os.Remove(destination)
	if err_ != nil || err != nil {
		return fmt.Errorf("%s %s", err, err_)
	} else if err != nil {
		return err
	} else {
		return err_
	}
}
