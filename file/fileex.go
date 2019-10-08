package file

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

//文件扩展类
type FileEx struct {
	os.File
	Path    string
	AbsPath string
	Handle  *os.File
}

//打开文件
func Open(path string, createIfNotExists, appendMode bool) (*FileEx, error) {
	fe := &FileEx{}
	_, err := fe.Open(path, createIfNotExists, appendMode)
	return fe, err
}

//获取绝对路径
func GetAbsPath(path string) string {
	absPath, _ := filepath.Abs(path)
	return absPath
}

//打开文件
//const (
//    O_RDONLY int = syscall.O_RDONLY // open the file read-only.
//    O_WRONLY int = syscall.O_WRONLY // open the file write-only.
//    O_RDWR   int = syscall.O_RDWR   // open the file read-write.
//    O_APPEND int = syscall.O_APPEND // 在文件末尾追加，打开后cursor在文件结尾位置
//    O_CREATE int = syscall.O_CREAT  // 如果不存在则创建
//    O_EXCL   int = syscall.O_EXCL   //与O_CREATE一起用，构成一个新建文件的功能，它要求文件必须不存在
//    O_SYNC   int = syscall.O_SYNC   // 同步方式打开，没有缓存，这样写入内容直接写入硬盘，系统掉电文件内容有一定保证
//    O_TRUNC  int = syscall.O_TRUNC  // 打开并清空文件
//)
//perm是文件的unix权限位，可以直接用数字写，如0644。可选值有：
//const (
//    // The single letters are the abbreviations
//    // used by the String method's formatting.
//    ModeDir        FileMode = 1 << (32 - 1 - iota) // d: is a directory
//    ModeAppend                                     // a: append-only
//    ModeExclusive                                  // l: exclusive use
//    ModeTemporary                                  // T: temporary file (not backed up)
//    ModeSymlink                                    // L: symbolic link
//    ModeDevice                                     // D: device file
//    ModeNamedPipe                                  // p: named pipe (FIFO)
//    ModeSocket                                     // S: Unix domain socket
//    ModeSetuid                                     // u: setuid
//    ModeSetgid                                     // g: setgid
//    ModeCharDevice                                 // c: Unix character device, when ModeDevice is set
//    ModeSticky
//
//    // Mask for the type bits. For regular files, none will be set.
//    ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice
//
//    ModePerm FileMode = 0777 // permission bits
//)
//上述可选值是权限位的高有效位，低有效位的值还是要用户自己写数字。
func (f *FileEx) openFile(path string, createIfNotExists, appendMode bool) (*os.File, error) {
	f.AbsPath = GetAbsPath(path)

	if !f.IsExist(f.AbsPath) {
		if !createIfNotExists {
			return nil, errors.New(fmt.Sprintf("[openFile] file《%s》does not exist", f.AbsPath))
		}

		file, err := os.Create(f.AbsPath)
		if file != nil {
			f.Path = path
			f.Handle = file
		}
		return file, err
	} else {
		mode := os.O_RDWR
		if appendMode {
			mode += os.O_APPEND
		} else {
			mode += os.O_TRUNC
		}

		file, err := os.OpenFile(f.AbsPath, mode, 0660)
		if file != nil {
			f.Path = path
			f.Handle = file
		}
		return file, err
	}
}

//打开文件
func (f *FileEx) Open(path string, createIfNotExists, appendMode bool) (*os.File, error) {
	if f.Handle != nil {
		absPath := GetAbsPath(path)

		if f.AbsPath == absPath {
			return f.Handle, nil
		} else {
			err := f.Close()
			if err != nil {
				return nil, err
			} else {
				return f.openFile(path, createIfNotExists, appendMode)
			}
		}
	} else {
		return f.openFile(path, createIfNotExists, appendMode)
	}
}

//关闭文件
func (f *FileEx) Close() error {
	if nil == f.Handle {
		return errors.New("文件句柄为空")
	}
	err := f.Handle.Close()
	if nil == err {
		f.Handle = nil
		f.Path = ""
		f.AbsPath = ""
	}
	return err
}

//删除文件
func (f *FileEx) Remove() error {
	if "" == f.AbsPath {
		return errors.New("文件路径为空")
	}

	f.Close()

	return os.Remove(f.AbsPath)
}

//获取文件信息
func (f *FileEx) Info() (os.FileInfo, error) {
	if "" == f.AbsPath {
		return nil, errors.New("文件路径为空")
	} else {
		return os.Stat(f.AbsPath)
	}
}

//获取文件的创建时间
func (f *FileEx) CreateTime() (int64, error) {
	fi, err := f.Info()
	if err != nil {
		return 0, err
	} else {
		return fi.ModTime().Unix(), nil
	}
}

//获取文件大小
func (f *FileEx) Size() (int64, error) {
	fi, err := f.Info()
	if err != nil {
		return 0, err
	} else {
		return fi.Size(), nil
	}
}

//重命名文件
func (f *FileEx) Rename(newPath string) error {
	newAbsPath := GetAbsPath(newPath)
	err := os.Rename(f.AbsPath, newAbsPath)
	if nil == err {
		f.Path = newPath
		f.AbsPath = newAbsPath
	}
	return err
}

//判断文件是否存在
func (f *FileEx) IsExist(path string) bool {
	_, err := os.Stat(GetAbsPath(path))
	return nil == err || os.IsExist(err)
}

//判断是否为文件
func (f *FileEx) IsFile(path string) bool {
	file, err := os.Stat(GetAbsPath(path))
	if err != nil {
		return false
	}
	return !file.IsDir()
}

//读取文件
func (f *FileEx) ReadEx(buff []byte, offset int64, unzipFile bool) (int, error) {
	if f.Handle != nil {
		if unzipFile {
			gr, err := f.newGzipReader()
			if err != nil {
				return 0, err
			}

			defer gr.Close()
			return gr.Read(buff)
		} else {
			if offset > 0 {
				return f.Handle.ReadAt(buff, offset)
			} else {
				return f.Handle.Read(buff)
			}
		}
	} else {
		return 0, errors.New("文件句柄为空")
	}
}

//将字符串写入文件
func (f *FileEx) WriteStringEx(content string) (int, error) {
	if f.Handle != nil {
		return f.Handle.WriteString(content)
	} else {
		return 0, errors.New("文件句柄为空")
	}
}

//写入文件：zipFile - 是否压缩文件
func (f *FileEx) WriteEx(data []byte, offset int64, zipFile bool) (int, error) {
	if f.Handle != nil {
		if zipFile {
			gw, err := f.newGzipWriter()
			if err != nil {
				return 0, err
			}
			defer gw.Close()
			return gw.Write(data[offset:])
		} else {
			if offset > 0 {
				return f.Handle.WriteAt(data, offset)
			} else {
				return f.Handle.Write(data)
			}
		}

	} else {
		return 0, errors.New("文件句柄为空")
	}
}

//创建Gzip压缩的读取流
func (f *FileEx) newGzipReader() (*gzip.Reader, error) {
	if nil == f.Handle {
		return nil, errors.New("文件句柄为空")
	} else {
		return gzip.NewReader(f.Handle)
	}
}

//创建Gzip压缩的写入流
func (f *FileEx) newGzipWriter() (*gzip.Writer, error) {
	if nil == f.Handle {
		return nil, errors.New("文件句柄为空")
	} else {
		gw := gzip.NewWriter(f.Handle)
		return gw, nil
	}
}

//创建Gzip.Tar压缩的读取流
func (f *FileEx) newGzipTarReader() (*gzip.Reader, *tar.Reader, error) {
	if nil == f.Handle {
		return nil, nil, errors.New("文件句柄为空")
	} else {
		gr, err := gzip.NewReader(f.Handle)
		if err != nil {
			return nil, nil, err
		} else {
			tr := tar.NewReader(gr)
			return gr, tr, nil
		}
	}
}

//创建Gzip.Tar压缩的写入流
func (f *FileEx) newGzipTarWriter() (*gzip.Writer, *tar.Writer, *tar.Header, error) {
	if nil == f.Handle {
		return nil, nil, nil, errors.New("文件句柄为空")
	} else {
		fi, err := f.Info()
		if err != nil {
			return nil, nil, nil, err
		} else {
			gw := gzip.NewWriter(f.Handle)
			tw := tar.NewWriter(gw)
			h := &tar.Header{Name: fi.Name(),
				Size:    fi.Size(),
				ModTime: fi.ModTime(),
				Mode:    int64(fi.Mode())}
			return gw, tw, h, nil
		}
	}
}
