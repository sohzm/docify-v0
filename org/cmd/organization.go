package main

import (
    "time"; "fmt"; "log"; "net/http";
)

func main() {
    fileServer := http.FileServer(http.Dir("./web"))
    http.Handle("/", fileServer)

    time_now := time.Now().Format(time.UnixDate)
    fmt.Printf("Organization Server\n" +
               "Time: " + time_now + "\n" +
               "Starting organization server at http://localhost:8087\n")

    if err := http.ListenAndServe(":8087", nil); err != nil {
        log.Fatal(err)
    }
}
