# Go Lockfile

Developed on Linux.

- creates a <filename>.lock file to mark that the file is locked for processing,

- function checks that the .lock file does not exist currently and attempts to lock it.

- triggers the callback where a custom code can run

- lastly, it will clean up the lock file.


```
package main

import (
	"github.com/rebooting/go_lockfile"
)

func main() {
	lf := go_lockfile.New("your-process-identifier-aaa-bbb", true)

	err := lf.LockRun("/tmp/nofile", func(f string) {
        // once file is locked, run your logic here
	})
	if err != nil {
		println(err.Error())
	}
}
```
