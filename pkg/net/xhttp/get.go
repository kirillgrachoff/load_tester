package xhttp

import (
    "net/http"
    "time"
)

type Response struct {
    Response *http.Response
    Err      error
    Time     time.Duration
}

func Get(url string) <-chan Response {
    c := make(chan Response)
    go func() {
        start := time.Now()
        resp, err := http.Get(url)
        c <- Response{resp, err, time.Since(start)}
        close(c)
    }()
    return c
}
