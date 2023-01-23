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

// Lockfile struct for creating .lock file
type LockFile struct {
}

// creates an empty LockFile
func New() LockFile {
	return LockFile{}
}

// creates the lockfile of the supplied filename
// it creates new file from fileName and a  ".lock" as the extention.
// the callback is executed after the lockfile is successfully created
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
