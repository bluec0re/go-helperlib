Go Library of my helper functions
=================================

Logging
-------

Example:
```go
import "github.com/bluec0re/go-helperlib/log"

func foo() {
	l := log.NewContextLogger()
	l.Infof("Log message in foo")
	l.Errorf("error message in foo")
}

func main() {
	_, err := log.AddFileHandler("test.log")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	log.Debugf("This is a %s msg", log.LevelDebug)
	log.Infof("This is a %s msg", log.LevelInfo)
	log.Warnf("This is a %s msg", log.LevelWarn)
	log.Errorf("This is a %s msg", log.LevelError)
	log.Fatalf("This is a %s msg", log.LevelFatal)
	foo()
}
```

Cli output:
<pre style="background-color:black; font-family:monospace">
[<span style="color:cyan">DEBUG</span>] This is a DEBUG msg
[<span style="color:#66f"> INFO</span>] This is a INFO msg
[<span style="color:yellow"> WARN</span>] This is a WARN msg
[<span style="color:red">ERROR</span>] This is a ERROR msg
[<span style="color:black; background-color:red">FATAL</span>] This is a FATAL msg
[<span style="color:#66f"> INFO</span>] Log message in foo
[<span style="color:red">ERROR</span>] error message in foo
</pre>



Logfile output:
```
[2018-03-12T13:07:33+01:00] DEBUG: <main (example.go:16)>: This is a DEBUG msg
[2018-03-12T13:07:33+01:00] INFO: <main (example.go:17)>: This is a INFO msg
[2018-03-12T13:07:33+01:00] WARN: <main (example.go:18)>: This is a WARN msg
[2018-03-12T13:07:33+01:00] ERROR: <main (example.go:19)>: This is a ERROR msg
[2018-03-12T13:07:33+01:00] FATAL: <main (example.go:20)>: This is a FATAL msg
[2018-03-12T13:07:33+01:00] f1ececd4-25ed-11e8-8411-a08cfde31c8b INFO: <foo (example.go:7)>: Log message in foo
[2018-03-12T13:07:33+01:00] f1ececd4-25ed-11e8-8411-a08cfde31c8b ERROR: <foo (example.go:8)>: error message in foo
```
