package go_lockfile_test

import (
	"github.com/rebooting/go_lockfile"
	"log"
	"os"
	"testing"
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

		lf := go_lockfile.New()
		if err := lf.Lock(eachTestCase.file, func(x string) {}); err != nil {
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
	lf := go_lockfile.New()
	tcase.fnSetup()

	lf.Lock(tcase.file, func(f string) {
		func() {
			t.Log(("waiting\n"))

			cf := go_lockfile.New()
			cerr := cf.Lock(tcase.file, func(f string) {
				t.Log("attempting to lock")
			})
			if cerr != nil {
				t.Logf("locking error %v", cerr)
			}else{
				t.Error("It should not lock successfully")
			}
	
			t.Log(("finished waiting\n"))
			defer tcase.fnTeardown()
		}()
	})
	
}
