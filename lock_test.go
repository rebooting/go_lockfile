package go_lockfile_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/rebooting/go_lockfile"
	// "time"
)

type testCase struct {
	file       string
	err        error
	fnSetup    func()
	fnLogic    func() error
	fnTeardown func()
}

// setup files for test cases
func setupAccess(t *testing.T, fileName string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func teardownAccess(t *testing.T, filename string) {
	t.Logf("removing. .. %s\n", filename)
	e := os.Remove(filename)
	if e != nil {
		t.Logf("trying to remove file in cleanup, error %v", e)
	}
}

func TestBasicLockCreation(t *testing.T) {

	testcases := []testCase{
		{
			file:    "/tmp/nofile",
			err:     go_lockfile.FileNotFoundError{},
			fnSetup: func() {},
			fnLogic: func() error { return nil },
			fnTeardown: func() {
				teardownAccess(t, "/tmp/nofile")
			},
		},
		{
			file:    "/tmp/successfile.log",
			err:     nil,
			fnSetup: func() { setupAccess(t, "/tmp/successfile.log") },
			fnLogic: func() error { return nil },
			fnTeardown: func() {
				teardownAccess(t, "/tmp/successfile.log")
			},
		},
	}

	for i, eachTestCase := range testcases {
		eachTestCase.fnSetup()
		defer eachTestCase.fnTeardown()

		lf := go_lockfile.New("aaa", go_lockfile.Options{Logging: true})
		if err := lf.LockRun(eachTestCase.file, func(x string) {}); err != nil {
			eachTestCase.fnLogic()
			if err != eachTestCase.err {
				t.Errorf("# %d Expected %v, got %v", i, eachTestCase.err, err.Error())
			}
		}
	}

}

func TestFileLocking(t *testing.T) {

	tcase := testCase{
		file:    "/tmp/nofile",
		err:     go_lockfile.FileIsLockedError{},
		fnSetup: func() { setupAccess(t, "/tmp/nofile") },
		fnLogic: func() error { return nil },
		fnTeardown: func() {
			teardownAccess(t, "/tmp/nofile")
		},
	}
	lf := go_lockfile.New("aaa", go_lockfile.Options{Logging: true})
	tcase.fnSetup()

	lf.LockRun(tcase.file, func(f string) {
		func() {
			t.Log(("waiting\n"))

			cf := go_lockfile.New("aaa", go_lockfile.Options{Logging: true})
			cerr := cf.LockRun(tcase.file, func(f string) {
				t.Log("attempting to lock")
			})
			if cerr != nil {
				t.Logf("locking error %v", cerr)
			} else {
				t.Error("It should not lock successfully")
			}

			t.Log(("finished waiting\n"))
			defer tcase.fnTeardown()
		}()
	})
}

func TestContentofLockfile(t *testing.T) {
	lf := go_lockfile.New("aaa-bbb", go_lockfile.Options{Logging: true})
	setupAccess(t, "/tmp/nofile")
	defer teardownAccess(t, "/tmp/nofile")
	lf.LockRun("/tmp/nofile", func(f string) {
		//linux locks are advisory
		lockfile := filepath.Clean(f + ".lock")
		data, err := os.ReadFile(lockfile)
		if err != nil {
			t.Errorf("can't read lock file %v", err.Error())
		}
		if string(data) != "aaa-bbb" {
			t.Errorf("expecting aaa-bbb got %s\n", data)
		}
	})
}
