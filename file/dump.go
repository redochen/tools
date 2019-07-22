package file

import (
	"fmt"
)

//保存到文件
func DumpFile(path, content string) error {
	fe, err := Open(path, true, false)
	if err != nil {
		fmt.Print(err)
		return err
	}

	defer fe.Close()

	_, err = fe.WriteStringEx(content)
	if err != nil {
		fmt.Print(err)
	}

	return err
}
