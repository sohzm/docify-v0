package main

import (
    "time"; "fmt"; "log"; "net/http";
)

func main() {
    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/", fileServer)

    time_now := time.Now().Format(time.UnixDate)
    fmt.Printf("Portal Server\n" +
               "Time: " + time_now + "\n" +
               "Starting portal server at http://localhost:8080\n")

    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
