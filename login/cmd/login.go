package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "time"
    "bufio"
    "context"
    "strconv"
    "net/http"
    "io/ioutil"
    "encoding/hex"
    "crypto/sha512"
    "path/filepath"
    "github.com/web3-storage/go-w3s-client"
)


func main() {
    fileServer := http.FileServer(http.Dir("./web"))
    http.Handle("/", fileServer)

    http.HandleFunc("/form", formHandler)
    http.HandleFunc("/list", listHandler)
    http.HandleFunc("/max_block", blockHandler)
    http.HandleFunc("/newaccount", accountHandler)

    time_now := time.Now().Format(time.UnixDate)
    fmt.Printf("Login Server\n" +
    "Time: " + time_now + "\n" +
    "Starting login server at http://localhost:8089\n")

    if err := http.ListenAndServe(":8089", nil); err != nil {
        log.Fatal(err)
    }
}

func listHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    fmt.Fprintf(w, "{\"links\" : [\"http://localhost:8087\", \"lalit\"]}")
}

func blockHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    b, err := ioutil.ReadFile("./blockchain/conf")
    if err != nil {
        fmt.Print(err)
    }
    block_num_str := string(b)
    fmt.Fprintf(w, "%s", block_num_str)
}

func formHandler(w http.ResponseWriter, r *http.Request) {

    // auth
    ad_username := r.FormValue("username")
    ad_password := r.FormValue("password")

    st_username := r.FormValue("std_user")

    if (ad_username != "admin" || ad_password != "pass") {
        fmt.Fprintf(w, "Admin username or password wrong")
        return
    }

    sha_512 := sha512.New()
    sha_512.Write([]byte("salt" + st_username))
    sha512_hash_user := fmt.Sprintf("%x", sha_512.Sum(nil))

    if _, err := os.Stat("user_data/" + sha512_hash_user); err != nil {
        fmt.Printf("%s %s", st_username, sha512_hash_user)
        fmt.Fprintf(w, "User doesnt exists")
        return;
    }
    fmt.Printf("Username Exists: %s", st_username)
    std_user := r.FormValue("std_user")
    //std_pass := r.FormValue("std_pass")
    doc_type := r.FormValue("doc_type")

    file, fileHeader, err := r.FormFile("file")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        fmt.Fprintf(w, "1")
        return
    }
    defer file.Close()

    // copy file to uploads directory
    time_i64 := time.Now().UTC().UnixNano()
    time_str := strconv.FormatInt(time_i64, 10)
    file_name := time_str + filepath.Ext(fileHeader.Filename)
    file_path := "uploads/" + file_name

    uploaded_file, err := os.Create(fmt.Sprintf("./uploads/%s", file_name))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Fprintf(w, "Error Uploading FIle to admin server\n")
        return
    }
    defer uploaded_file.Close()

    _, err = io.Copy(uploaded_file, file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Fprintf(w, "4")
        return
    }
    fmt.Fprintf(w, "Upload on admin server successful...\n")

    // Upload on IPFS Network
    w3s_context, err := w3s.NewClient(w3s.WithToken(os.Getenv("Token")))
    if err != nil {
        fmt.Printf("")
        fmt.Fprintf(w, "<h1>File not uploaded to IPFS\n</h1><a href='http://localhost:8089/form.html'>Go to Admin Panel</a> ")
        return
    }
    f, err := os.Open(file_path)
    if err != nil {
        fmt.Fprintf(w, "<h1>File not uploaded to IPFS\n</h1><a href='http://localhost:8089/form.html'>Go to Admin Panel</a> ")
        return
    }
    defer f.Close()

    cid, _ := w3s_context.Put(context.Background(), f)

    fmt.Fprintf(w, fmt.Sprintf("https://%v.ipfs.dweb.link\n", cid))
    //
    b, err := ioutil.ReadFile("./blockchain/conf")
    if err != nil {
        fmt.Print(err)
    }
    block_num_str := string(b)
    block_num_int, _ := strconv.Atoi(block_num_str)

    ipfs_link_str := fmt.Sprintf("https://%v.ipfs.dweb.link\n", cid)

    ux, err := ioutil.ReadFile("./blockchain/" + block_num_str)
    if err != nil {
        fmt.Print(err)
    }

    block_num_int++
    sp := string(ux)
    hp := sha512.New()
    hp.Write([]byte(sp))
    sha512_hash := hex.EncodeToString(hp.Sum(nil))

    histr := strconv.Itoa(block_num_int)

    block_str := "{ \"block\": \"" + histr + "\", \"" + std_user + "\": [\"" + ipfs_link_str[:len(ipfs_link_str) - 1] + "\", \"" + doc_type + "\"], \"timestamp\": \"" + time_str + "\", \"prev\": \"" + sha512_hash + "\" }"
    fmt.Printf(block_str + "\n")

    f45, err := os.Create("blockchain/" + histr)
    if err != nil {
        log.Fatal(err)
    }
    defer f45.Close()
    _, err2 := f45.WriteString(block_str)
    if err2 != nil {
        log.Fatal(err2)
    }

    fx1, _ := os.Create("blockchain/conf")
    defer f.Close()

    // adding
    bytesRead, err := ioutil.ReadFile("user_data/" + sha512_hash_user) 
    s := string(bytesRead)
    neh := s[:16]  + "\n        [\"" + doc_type + "\", " + histr + ", \"" + sha512_hash + "\", \"" + time_str + "\"]," + s[16:];
    ixf, _ := os.Create("user_data/" + sha512_hash_user)
    defer ixf.Close()
    ixf.WriteString(neh)


    wx1 := bufio.NewWriter(fx1)
    fmt.Fprintf(wx1, histr)
    wx1.Flush()
    fmt.Fprintf(w, "<h1> File Successfully Uploaded </h1><a href='http://localhost:8089/form.html'>Go to Admin Panel</a> ")
}

func generalHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "connection successful")
}

func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
}


func accountHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    // auth
    ad_username := r.FormValue("username")
    ad_password := r.FormValue("password")

    st_username := r.FormValue("std_user")

    if (ad_username != "admin" || ad_password != "pass") {
        fmt.Fprintf(w, "Admin username or password wrong")
        return
    }
    sha_512 := sha512.New()
    sha_512.Write([]byte("salt" + st_username))
    sha512_hash_user := fmt.Sprintf("%x", sha_512.Sum(nil))
    src := "user_data/temp"
    dest := "user_data/" + sha512_hash_user

    if _, err := os.Stat("user_data/" + sha512_hash_user); err == nil {
        fmt.Printf("%s %s", st_username, sha512_hash_user)
        fmt.Fprintf(w, "<h1> Account already exists</h1><a href='http://localhost:8089/form.html'>Go to Admin Panel</a> ")
        return;
    }


    bytesRead, err := ioutil.ReadFile(src)

    if err != nil {
        fmt.Fprintf(w, "Error while creating account, try again")
        return

    }

    err = ioutil.WriteFile(dest, bytesRead, 0644)

    if err != nil {
        fmt.Fprintf(w, "Error while creating account, try again")
        return
    }
    fmt.Fprintf(w, "<h1> Account created </h1><a href='http://localhost:8089/form.html'>Go to Admin Panel</a> ")
}
          
