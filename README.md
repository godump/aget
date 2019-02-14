# Package aget

Package aget provides functions for opens the named file for reading.

```go
import "github.com/godump/aget"
```

- [Examples](#Examples)

# Examples

`Open` select the appropriate method to open the file based on the incoming args automatically:

```go
rc, err := aget.Open("/etc/hosts")
...
rc, err := aget.Open("https://github.com/godump/aget/blob/master/README.md")
...
```

Must call `Close()` when finished with it:

```go
rc.Close()
```

Using local disk storage to cache content:

```go
rc, err := aget.OpenEx("https://github.com/godump/aget/blob/master/README.md", "/tmp/README.md", time.Hour)
...
```
