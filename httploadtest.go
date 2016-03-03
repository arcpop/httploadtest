package main

import (
    "net"
    "strings"
    "sync"
    "time"
    "fmt"
)

func contains(b []byte, s string) bool {
    for i := 0; i < len(b) - len(s); i++ {
        if b[i] == s[0] {
            equal := true
            for j := 1; j < len(s); j++ {
                if b[i + j] != s[j] {
                    equal = false
                    break
                }
            }
            if equal {
                return true
            }
        }
    }
    return false
}

type stringError struct {
    what string
}

func (e *stringError) Error() string {
    return e.what
}

func newError(s string) *stringError {
    return &stringError{ what: s }
}

func runClient()  (e error) {
    c, e := net.Dial("tcp", "debian-server:80")
    
    if e != nil {
        return 
    }
    r := []byte("")
    n := 0
    for n < len(r) {
        var i int
        i, e = c.Write(r[n:])
        if e != nil {
            return
        }
        n += i
    }
    
    n = 0
    resp := make([]byte, 0, 1024)
    for !contains(resp, "\r\n") {
        i := 0
        i, e = c.Read(resp[n:])
        if e != nil {
            return
        }
        n += i
    }
    
    ans := string(resp[:n])
    
    if strings.HasPrefix(ans, "HTTP/1.1 200") {
        return nil
    }
    
    return newError("Bad HTTP answer")
}

type resultHolder struct {
    m sync.Mutex
    oks int
    errors int
}

var result resultHolder

func run()  {
    for {
        e := runClient()
        result.m.Lock()
        if e == nil {
            result.oks++
        } else {
            fmt.Println(e)
            result.errors++;
        }
        result.m.Unlock()
    }
}

func runClients() {
    for i := 0; i < 1; i++ {
        go run()
    }
}

func main()  {
    go runClients()
    next := time.Now().Add(1 * time.Second)
    for {
        time.Sleep(next.Sub(time.Now()))
        result.m.Lock()
        o := result.oks
        e := result.errors
        result.oks = 0
        result.errors = 0
        result.m.Unlock()
        fmt.Println("OKs: ", o, " | Errors: ", e)
        next = next.Add(1 * time.Second)
    }
}