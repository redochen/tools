package file

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

//Helper 文件帮助类
type Helper struct {
	os.File
	Path    string
	AbsPath string
	Handle  *os.File
}

//Open 打开文件
func Open(path string, createIfNotExists, appendMode bool) (*Helper, error) {
	fe := &Helper{}
	_, err := fe.Open(path, createIfNotExists, appendMode)
	return fe, err
}

//GetAbsPath 获取绝对路径
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
//openFile 上述可选值是权限位的高有效位，低有效位的值还是要用户自己写数字。
func (h *Helper) openFile(path string, createIfNotExists, appendMode bool) (*os.File, error) {
	h.AbsPath = GetAbsPath(path)

	if !h.IsExist(h.AbsPath) {
		if !createIfNotExists {
			return nil, fmt.Errorf("[openFile] file《%s》does not exist", h.AbsPath)
		}

		file, err := os.Create(h.AbsPath)
		if file != nil {
			h.Path = path
			h.Handle = file
		}

		return file, err
	}

	mode := os.O_RDWR
	if appendMode {
		mode += os.O_APPEND
	} else {
		mode += os.O_TRUNC
	}

	file, err := os.OpenFile(h.AbsPath, mode, 0660)
	if file != nil {
		h.Path = path
		h.Handle = file
	}

	return file, err
}

//Open 打开文件
func (h *Helper) Open(path string, createIfNotExists, appendMode bool) (*os.File, error) {
	if nil == h.Handle {
		return h.openFile(path, createIfNotExists, appendMode)
	}

	absPath := GetAbsPath(path)
	if h.AbsPath == absPath {
		return h.Handle, nil
	}

	err := h.Close()
	if err != nil {
		return nil, err
	}

	return h.openFile(path, createIfNotExists, appendMode)
}

//Close 关闭文件
func (h *Helper) Close() error {
	if nil == h.Handle {
		return errors.New("文件句柄为空")
	}

	err := h.Handle.Close()
	if nil == err {
		h.Handle = nil
		h.Path = ""
		h.AbsPath = ""
	}

	return err
}

//Remove 删除文件
func (h *Helper) Remove() error {
	if "" == h.AbsPath {
		return errors.New("文件路径为空")
	}

	h.Close()

	return os.Remove(h.AbsPath)
}

//Info 获取文件信息
func (h *Helper) Info() (os.FileInfo, error) {
	if "" == h.AbsPath {
		return nil, errors.New("文件路径为空")
	}

	return os.Stat(h.AbsPath)
}

//CreateTime 获取文件的创建时间
func (h *Helper) CreateTime() (int64, error) {
	fi, err := h.Info()
	if err != nil {
		return 0, err
	}

	return fi.ModTime().Unix(), nil
}

//Size 获取文件大小
func (h *Helper) Size() (int64, error) {
	fi, err := h.Info()
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

//Rename 重命名文件
func (h *Helper) Rename(newPath string) error {
	newAbsPath := GetAbsPath(newPath)
	err := os.Rename(h.AbsPath, newAbsPath)
	if nil == err {
		h.Path = newPath
		h.AbsPath = newAbsPath
	}

	return err
}

//IsExist 判断文件是否存在
func (h *Helper) IsExist(path string) bool {
	_, err := os.Stat(GetAbsPath(path))
	return nil == err || os.IsExist(err)
}

//IsFile 判断是否为文件
func (h *Helper) IsFile(path string) bool {
	file, err := os.Stat(GetAbsPath(path))
	if err != nil {
		return false
	}

	return !file.IsDir()
}

//ReadEx 读取文件
func (h *Helper) ReadEx(buff []byte, offset int64, unzipFile bool) (int, error) {
	if nil == h.Handle {
		return 0, errors.New("文件句柄为空")
	}

	if unzipFile {
		gr, err := h.newGzipReader()
		if err != nil {
			return 0, err
		}

		defer gr.Close()
		return gr.Read(buff)
	}

	if offset > 0 {
		return h.Handle.ReadAt(buff, offset)
	}

	return h.Handle.Read(buff)
}

//WriteStringEx 将字符串写入文件
func (h *Helper) WriteStringEx(content string) (int, error) {
	if nil == h.Handle {
		return 0, errors.New("文件句柄为空")
	}

	return h.Handle.WriteString(content)
}

//WriteEx 写入文件：zipFile - 是否压缩文件
func (h *Helper) WriteEx(data []byte, offset int64, zipFile bool) (int, error) {
	if nil == h.Handle {
		return 0, errors.New("文件句柄为空")
	}

	if zipFile {
		gw, err := h.newGzipWriter()
		if err != nil {
			return 0, err
		}

		defer gw.Close()
		return gw.Write(data[offset:])
	}

	if offset > 0 {
		return h.Handle.WriteAt(data, offset)
	}

	return h.Handle.Write(data)
}

//创建Gzip压缩的读取流
func (h *Helper) newGzipReader() (*gzip.Reader, error) {
	if nil == h.Handle {
		return nil, errors.New("文件句柄为空")
	}

	return gzip.NewReader(h.Handle)
}

//创建Gzip压缩的写入流
func (h *Helper) newGzipWriter() (*gzip.Writer, error) {
	if nil == h.Handle {
		return nil, errors.New("文件句柄为空")
	}

	gw := gzip.NewWriter(h.Handle)
	return gw, nil
}

//创建Gzip.Tar压缩的读取流
func (h *Helper) newGzipTarReader() (*gzip.Reader, *tar.Reader, error) {
	if nil == h.Handle {
		return nil, nil, errors.New("文件句柄为空")
	}

	gr, err := gzip.NewReader(h.Handle)
	if err != nil {
		return nil, nil, err
	}

	tr := tar.NewReader(gr)
	return gr, tr, nil
}

//创建Gzip.Tar压缩的写入流
func (h *Helper) newGzipTarWriter() (*gzip.Writer, *tar.Writer, *tar.Header, error) {
	if nil == h.Handle {
		return nil, nil, nil, errors.New("文件句柄为空")
	}

	fi, err := h.Info()
	if err != nil {
		return nil, nil, nil, err
	}

	gw := gzip.NewWriter(h.Handle)
	tw := tar.NewWriter(gw)
	head := &tar.Header{Name: fi.Name(),
		Size:    fi.Size(),
		ModTime: fi.ModTime(),
		Mode:    int64(fi.Mode())}

	return gw, tw, head, nil
}
