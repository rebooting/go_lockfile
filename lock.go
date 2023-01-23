package golockfile

import (
	"errors"
	"log"
	"os"
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

type LockFile struct {
}

func New() LockFile {
	return LockFile{}
}

func (l *LockFile) Lock(fileName string, callback func(string)) error {
	// check if file exist
	if _, err := os.Stat(fileName); err != nil {
		// log.Printf("cannot STAT file %s , error: %v\n", fileName, err)
		if errors.Is(err, os.ErrNotExist) {
			return FileNotFoundError{}
		}
	}
	// create lockfile
	file, err := os.OpenFile(fileName+".lock", os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	defer file.Close()
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
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	

	callback(fileName)

	return nil
}

func openFile(fileName string) (file *os.File, err error) {

	file, err = os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	return
}
