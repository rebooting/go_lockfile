package go_lockfile

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

type FileNotFoundError struct{}

func (m FileNotFoundError) Error() string {
	return "File Not Found"
}

type FileIsLockedError struct{}

func (m FileIsLockedError) Error() string {
	return "File Already Locked"
}

type TryLaterErr struct{}

func (m TryLaterErr) Error() string {
	return "File Locking in Action"
}

// Lockfile struct for creating .lock file
type LockFile struct {
	logging bool
	id      string
	mutex   sync.Mutex
}

// creates a LockFile with the caller ID andoption to turn on logging,
// the caller ID should be a unique identifier from the proces:
// e.g. PID, SNO or UUID in case if the lockfile is not deleted, the content of the lock file can be used to indicate if the process still exists before being deleted.
func New(id string, logging bool) LockFile {
	return LockFile{id: id, logging: logging}
}

// creates the lockfile of the supplied filename
// it creates new file from fileName and a  ".lock" as the extention.
// the callback is executed after the lockfile is successfully created
func (l *LockFile) LockRun(filePath string, runnableCallback func(string)) error {
	fileName := filepath.Join("", filepath.Clean(filePath))
	if l.logging {
		log.Printf("attempting to acquire mutex for %s\n", fileName)
	}

	l.mutex.Lock()
	if l.logging {
		log.Printf("mutex acquired for %s\n", fileName)
	}
	defer l.mutex.Unlock()
	// check if file exist
	if _, err := os.Stat(fileName); err != nil {
		if l.logging {
			log.Printf("cannot STAT file %s , error: %v\n", fileName, err)
		}
		if errors.Is(err, os.ErrNotExist) {
			return FileNotFoundError{}
		}
	}
	// create lockfile
	lockfile := filepath.Clean(filePath + ".lock")
	file, err := os.OpenFile(lockfile, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %s\n", err)
		}
	}()

	defer os.Remove(fileName + ".lock")
	if err != nil {
		log.Printf("cannot O_EXCL file %s , error: %v\n", fileName, err)
		return err
	}

	//get exclusive lock
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		log.Printf("cannot FLOCK file %s , error: %v\n", fileName, err)
		return FileIsLockedError{}
	}
	// write identifer to file
	n, err:=file.WriteString(l.id)
	if err!=nil{
		return err
	}
	if l.logging{
		log.Printf("%d bytes written to %s", n, lockfile)
	}
	
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

	runnableCallback(fileName)

	return nil
}
