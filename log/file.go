package log

import (
	"encoding/json"
	"fmt"
	. "github.com/redochen/tools/config"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	logFileOpton = "logpath"
	LogFileDir   = "log"
)

// FileLogWriter implements LoggerInterface.
// It writes messages by lines limit, file size limit, or time frequency.
type FileLogWriter struct {
	*log.Logger
	mw *MuxWriter
	// The opened file
	Filename string `json:"filename"`

	Maxlines          int `json:"maxlines"`
	maxlines_curlines int

	// Rotate at size
	Maxsize         int `json:"maxsize"`
	maxsize_cursize int

	// Rotate daily
	Daily          bool  `json:"daily"`
	Maxdays        int64 `json:"maxdays"`
	daily_opendate string

	Rotate bool `json:"rotate"`

	startLock sync.Mutex // Only one log can write to the file

	Level int `json:"level"`
}

// an *os.File writer with locker.
type MuxWriter struct {
	sync.Mutex
	fd *os.File
}

// write to os.File.
func (l *MuxWriter) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.fd.Write(b)
}

// set os.File in writer.
func (l *MuxWriter) SetFd(fd *os.File) {
	if l.fd != nil {
		l.fd.Close()
	}
	l.fd = fd
}

// create a FileLogWriter returning as LoggerInterface.
func NewFileWriter() LoggerInterface {
	w := &FileLogWriter{
		Filename: "",
		Maxlines: 1000000,
		Maxsize:  1 << 28, //256 MB
		Daily:    true,
		Maxdays:  7,
		Rotate:   true,
		Level:    LevelTrace,
	}
	// use MuxWriter instead direct use os.File for lock write when rotate
	w.mw = new(MuxWriter)
	// set MuxWriter as Logger's io.Writer
	w.Logger = log.New(w.mw, "", 0)
	return w
}

// Init file logger with json config.
// jsonconfig like:
//	{
//	"filename":"logs/beego.log",
//	"maxlines":10000,
//	"maxsize":1<<30,
//	"daily":true,
//	"maxdays":15,
//	"rotate":true
//	}
func (w *FileLogWriter) Init(jsonconfig string) error {
	err := json.Unmarshal([]byte(jsonconfig), w)
	if err != nil {
		return err
	}

	/*
		if len(w.Filename) == 0 {
			return errors.New("jsonconfig must have filename")
		}
	*/
	_, err = os.Stat(LogFileDir)
	if err != nil {
		os.Mkdir(LogFileDir, os.ModeDir)
	}
	now := time.Now()
	folderName := now.Format("2006-01-02")
	_, err = os.Stat(fmt.Sprintf("%s/%s", LogFileDir, folderName))
	if err != nil {
		os.Mkdir(fmt.Sprintf("%s/%s", LogFileDir, folderName), os.ModeDir)
	}
	w.Filename = fmt.Sprintf("%s/%s/%s.log", LogFileDir, folderName, now.Format("2006-01-02.15"))
	err = w.startLogger()
	return err
}

// start file logger. create log file and set to locker-inside file writer.
func (w *FileLogWriter) startLogger() error {
	fd, err := w.createLogFile()
	if err != nil {
		return err
	}
	w.mw.SetFd(fd)
	err = w.initFd()
	if err != nil {
		return err
	}
	return nil
}

func (w *FileLogWriter) docheck(size int) {
	w.startLock.Lock()
	defer w.startLock.Unlock()
	if w.Rotate && ((w.Maxlines > 0 && w.maxlines_curlines >= w.Maxlines) ||
		(w.Maxsize > 0 && w.maxsize_cursize >= w.Maxsize)) {
		if err := w.Rename(); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
			return
		}
	} else if w.Rotate && w.Daily && time.Now().Format("2006010215") != w.daily_opendate {
		if err := w.DoRotate(); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
			return
		}
	}
	w.maxlines_curlines++
	w.maxsize_cursize += size
}

// write logger message into file.
func (w *FileLogWriter) WriteMsg(msg string, level int) error {
	if level > w.Level || level == LevelWarning {
		return nil
	}
	message := fmt.Sprintf(`{"level":"%s", "time":"%s", "machine":"%s", "message":"%s"}`,
		getLevelName(level),
		time.Now().Local().Format("2006-01-02 15:04:05"),
		machineName,
		msg)
	n := 24 + len(message) // 24 stand for the length "2013/06/23 21:00:22 [T] "
	w.docheck(n)
	w.Logger.Println(message)
	return nil
}

func (w *FileLogWriter) createLogFile() (*os.File, error) {
	// Open the log file

	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	return fd, err
}

func (w *FileLogWriter) initFd() error {
	fd := w.mw.fd
	finfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s\n", err)
	}
	w.maxsize_cursize = int(finfo.Size())
	w.daily_opendate = time.Now().Format("2006010215")
	if finfo.Size() > 0 {
		content, err := ioutil.ReadFile(w.Filename)
		if err != nil {
			return err
		}
		w.maxlines_curlines = len(strings.Split(string(content), "\n"))
	} else {
		w.maxlines_curlines = 0
	}
	return nil
}

//rename log file
func (w *FileLogWriter) Rename() error {
	_, err := os.Lstat(w.Filename)
	if err != nil {
		//file not exists
		return nil
	}

	// Find the next available number
	num := 1
	fname := ""
	for ; err == nil && num <= 999; num++ {
		fname = fmt.Sprintf("%s.%03d", w.Filename, num)
		_, err = os.Lstat(fname)
	}
	// return error if the last file checked still existed
	if err == nil {
		return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", w.Filename)
	}

	// block Logger's io.Writer
	w.mw.Lock()
	defer w.mw.Unlock()
	w.mw.fd.Close()

	// close fd before rename
	// Rename the file to its newfound home
	err = os.Rename(w.Filename, fname)
	if err != nil {
		return fmt.Errorf("rename: %s\n", err)
	}

	now := time.Now()
	folderName := now.Format("2006-01-02")
	_, err = os.Stat(fmt.Sprintf("%s/%s", LogFileDir, folderName))
	if err != nil {
		os.Mkdir(fmt.Sprintf("%s/%s", LogFileDir, folderName), os.ModeDir)
	}

	w.Filename = fmt.Sprintf("%s/%s/%s.log", LogFileDir, folderName, now.Format("2006-01-02.15"))
	w.daily_opendate = time.Now().Format("2006010215")
	// re-start logger
	err = w.startLogger()
	if err != nil {
		return fmt.Errorf("rename StartLogger: %s\n", err)
	}

	go w.deleteOldLog()
	return nil
}

// DoRotate means it need to write file in new file.
// new file name like xx.log.2013-01-01.2
func (w *FileLogWriter) DoRotate() error {
	now := time.Now()
	folderName := now.Format("2006-01-02")
	_, err := os.Stat(fmt.Sprintf("%s/%s", LogFileDir, folderName))
	if err != nil {
		os.Mkdir(fmt.Sprintf("%s/%s", LogFileDir, folderName), os.ModeDir)
	}

	w.Filename = fmt.Sprintf("%s/%s/%s.log", LogFileDir, folderName, now.Format("2006-01-02.15"))

	w.mw.Lock()
	defer w.mw.Unlock()

	w.mw.fd.Close()
	w.daily_opendate = time.Now().Format("2006010215")

	// re-start logger
	err = w.startLogger()
	if err != nil {
		return fmt.Errorf("Rotate StartLogger: %s\n", err)
	}
	go w.deleteOldLog()
	return nil
}

func (w *FileLogWriter) deleteOldLog() {
	dir := filepath.Dir(fmt.Sprintf("%s/%s", LogFileDir, time.Now().Format("2006-01-02")))
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				returnErr = fmt.Errorf("unable to delete old log '%s', error: %+v", path, r)
				fmt.Println(returnErr)
			}
		}()
		if info.ModTime().Unix() < (time.Now().Unix() - 60*60*24*w.Maxdays) {
			if info.IsDir() {
				if info.Name() != LogFileDir {
					s, _ := ioutil.ReadDir(path)
					if len(s) == 0 {
						os.RemoveAll(path)
					}
				}

			} else {
				os.Remove(path)
			}

		}
		return
	})
}

// destroy file logger, close file writer.
func (w *FileLogWriter) Destroy() {
	w.mw.fd.Close()
}

// flush file logger.
// there are no buffering messages in file logger in memory.
// flush file means sync file from disk.
func (w *FileLogWriter) Flush() {
	w.mw.fd.Sync()
}

func registerFile() {
	if nil == Conf || !Conf.IsValid() {
		return
	}

	if Conf != nil {
		value, err := Conf.String(LogSection, logFileOpton)
		if err != nil {
			fmt.Println(err)
		} else {
			LogFileDir = value
		}
	}

	Register("file", NewFileWriter)
}
