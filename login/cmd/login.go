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
    "crypto/sha1"
    "encoding/hex"
    "path/filepath"
    "github.com/web3-storage/go-w3s-client"
)


func main() {
    fileServer := http.FileServer(http.Dir("./web"))
    http.Handle("/", fileServer)

    http.HandleFunc("/form", formHandler)
    http.HandleFunc("/list", listHandler)

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

func formHandler(w http.ResponseWriter, r *http.Request) {

    // auth
    ad_username := r.FormValue("username")
    ad_password := r.FormValue("password")

    if (ad_username != "admin" || ad_password != "pass") {
        fmt.Fprintf(w, "Admin username or password wrong")
        return
    }

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
	    fmt.Fprintf(w, "3")
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
    w3s_context, _ := w3s.NewClient(w3s.WithToken(os.Getenv("Token")))
    f, err := os.Open(file_path)
    if err != nil {
        fmt.Printf("Error opening file\n")
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
    hp := sha1.New()
    hp.Write([]byte(sp))
    sha1_hash := hex.EncodeToString(hp.Sum(nil))

    histr := strconv.Itoa(block_num_int)

    block_str := "{ 'block': '" + histr + "', '" + std_user + "': ['" + ipfs_link_str + "', '" + doc_type + "'], 'timestamp': '" + time_str + "', 'prev': '" + sha1_hash + "' }"
    fmt.Printf(block_str)

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

    wx1 := bufio.NewWriter(fx1)
    fmt.Fprintf(wx1, histr)
    wx1.Flush()
}

func generalHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "connection successful")
}

func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
}
