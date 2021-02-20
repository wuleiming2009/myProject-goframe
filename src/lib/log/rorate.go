package log

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
	host     = "locahost"
	userName = "root"
	zlogDir  string
	zlogDirs []string
)

func init() {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}

	userName = strings.Replace(userName, `\`, "_", -1)

}

// If non-empty, overrides the choice of directory in which to write logs.
// See createLogDirs for the full list of possible destinations.
func createZapLogDirs() {
	if zlogDir != "" {
		zlogDirs = append(zlogDirs, zlogDir)
	}
	zlogDirs = append(zlogDirs, os.TempDir())
}

var MaxSize uint64 = 1024 * 1024 * 1800

type severity int32 // sync/atomic int32

// These constants identify the log levels in order of increasing severity.
// A message written to a high-severity log file is also written to each
// lower-severity log file.
const (
	debugLog severity = iota
	infoLog
	warningLog
	errorLog
	fatalLog
	numSeverity = 4
)

var onceLogDirs sync.Once
var flushInterval = time.Second * 10

// logDirs lists the candidate directories for new log files.

var severityName = []string{
	debugLog:   "DEBUG",
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}

type syncBuffer struct {
	mu sync.Mutex
	*bufio.Writer
	file   *os.File
	sev    severity
	nbytes uint64 // The number of bytes written to this file
}

func (sb *syncBuffer) Sync() error {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	sb.Flush()
	return sb.file.Sync()
}

func (sb *syncBuffer) flush() {
	sb.mu.Lock()
	sb.Flush()
	sb.mu.Unlock()
}

func (sb *syncBuffer) FlushDaemon() {
	for range time.NewTicker(flushInterval).C {
		sb.flush()
	}
}

func (sb *syncBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if sb.nbytes+uint64(len(p)) >= MaxSize {
		if err := sb.rotateFile(time.Now()); err != nil {
			return 0, err
		}
	}

	n, err = sb.Writer.Write(p)
	sb.nbytes += uint64(n)
	if err != nil {
		return 0, err
	}
	return
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (sb *syncBuffer) rotateFile(now time.Time) error {
	if sb.file != nil {
		err := sb.Flush()
		if err != nil {
			fmt.Println(err.Error())
		}
		err = sb.file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	var err error
	sb.file, _, err = create(now)
	sb.nbytes = 0
	if err != nil {
		return err
	}

	sb.Writer = bufio.NewWriterSize(sb.file, bufferSize)

	// Write header.
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "ZapLog file created at: %s\n", now.Format("2006/01/02 15:04:05"))
	fmt.Fprintf(&buf, "Running on machine: %s\n", host)
	fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	n, err := sb.file.Write(buf.Bytes())
	sb.nbytes += uint64(n)
	return err
}

// bufferSize sizes the buffer associated with each log file. It's large
// so that log records can accumulate without the logging thread blocking
// on disk I/O. The flushDaemon will block instead.
const bufferSize = 256 * 1024

func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

func logName(t time.Time) (name, link string) {
	name = fmt.Sprintf("zap.%s.%s.%s.%04d%02d%02d-%02d%02d%02d.%d.log",
		program,
		host,
		userName,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		pid)
	return name, fmt.Sprintf("%s.log", program)
}

func create(t time.Time) (f *os.File, filename string, err error) {
	onceLogDirs.Do(createZapLogDirs)
	if len(zlogDirs) == 0 {
		return nil, "", errors.New("log: no log dirs")
	}

	name, link := logName(t)
	var lastErr error
	for _, dir := range zlogDirs {
		fname := filepath.Join(dir, name)
		f, err := os.Create(fname)
		if err == nil {
			symlink := filepath.Join(dir, link)
			_ = os.Remove(symlink)        // ignore err
			_ = os.Symlink(name, symlink) // ignore err
			return f, fname, nil
		}
		lastErr = err
	}
	return nil, "", fmt.Errorf("log: cannot create log: %v", lastErr)
}
