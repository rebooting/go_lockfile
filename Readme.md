# Go Lockfile

Developed on Linux.



- creates a <filename>.lock file to mark that the file is locked for processing,

- function checks that the .lock file does not exist currently and attempts to lock it.

- triggers the callback where a custom code can run if there are no errors.

- lastly, it will clean up the lock file.


### To create a lock file that exists

to create a lock file for existing file 

e.g. /tmp/nofile

the library will check for the existence of /tmp/nofile before creating /tmp/nofile.lock


```
package main

import (
	"github.com/rebooting/go_lockfile"
)

func main() {
	lf := go_lockfile.New("your-process-identifier-aaa-bbb", go_lockfile.Options{})

	err := lf.LockRun("/tmp/nofile", func(f string) {
        // once file is locked, run your logic here
				// if there are errors, it will not execute
	})
	if err != nil {
		println(err.Error())
	}
}
```
### To just create a lock file without checking for the file name that exists

```
package main

import (
	"github.com/rebooting/go_lockfile"
)

func main() {
	lf := go_lockfile.New("your-process-identifier-aaa-bbb", go_lockfile.Options{NoFileDependency:true})

	err := lf.LockRun("/tmp/nofile", func(f string) {
        // once file is locked, run your logic here
				// if there are errors, it will not execute
	})
	if err != nil {
		println(err.Error())
	}
}
```


