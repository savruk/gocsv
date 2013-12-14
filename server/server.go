package server

import (
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/render"
    "net/http"
    "net/url"
    "fmt"
    "regexp"
    "mime/multipart"
    "log"
    "io"
    "io/ioutil"
    "bytes"
    "os"
    "path/filepath"
    "encoding/json"
    "strings"
)

type CsvServer struct{
    Port int
    Static string
    Media string
    Templates string
}

const (
    WEBSITE           = "http://blueimp.github.io/jQuery-File-Upload/"
    MIN_FILE_SIZE     = 1       // bytes
    MAX_FILE_SIZE     = 5000000000000 // bytes
    IMAGE_TYPES       = "image/(gif|p?jpeg|(x-)?png)"
    ACCEPT_FILE_TYPES = "text/csv"
    EXPIRATION_TIME   = 300 // seconds
    THUMBNAIL_PARAM   = "=s80"
)

var (
    imageTypes      = regexp.MustCompile(IMAGE_TYPES)
    acceptFileTypes = regexp.MustCompile(ACCEPT_FILE_TYPES)
)

type FileInfo struct {
    Url          string            `json:"url,omitempty"`
    ThumbnailUrl string            `json:"thumbnailUrl,omitempty"`
    Name         string            `json:"name"`
    Type         string            `json:"type"`
    Size         int64             `json:"size"`
    Error        string            `json:"error,omitempty"`
    DeleteUrl    string            `json:"deleteUrl,omitempty"`
    DeleteType   string            `json:"deleteType,omitempty"`
}

func (fi *FileInfo) ValidateType() (valid bool) {
    if acceptFileTypes.MatchString(fi.Type) {
        return true
    }
    fi.Error = "Filetype not allowed"
    return false
}

func (fi *FileInfo) ValidateSize() (valid bool) {
    if fi.Size < MIN_FILE_SIZE {
        fi.Error = "File is too small"
    } else if fi.Size > MAX_FILE_SIZE {
        fi.Error = "File is too big"
    } else {
        return true
    }
    return false
}

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func (cs *CsvServer) handleUpload(r *http.Request, p *multipart.Part) (fi *FileInfo) {
    fi = &FileInfo{
        Name: p.FileName(),
        Type: p.Header.Get("Content-Type"),
    }
    // if !fi.ValidateType() {
    //     return
    // }
    defer func() {
        if rec := recover(); rec != nil {
            log.Println(rec)
            fi.Error = rec.(error).Error()
        }
    }()
    lr := &io.LimitedReader{R: p, N: MAX_FILE_SIZE + 1}
    fo, _ := os.Create(filepath.Join(cs.Media, fi.Name))
    _, _ = io.Copy(fo, lr)
    return
}

func getFormValue(p *multipart.Part) string {
    var b bytes.Buffer
    io.CopyN(&b, p, int64(1<<20)) // Copy max: 1 MiB
    return b.String()
}

func (cs *CsvServer) handleUploads(r *http.Request) (fileInfos []*FileInfo) {
    fileInfos = make([]*FileInfo, 0)
    mr, err := r.MultipartReader()
    check(err)
    r.Form, err = url.ParseQuery(r.URL.RawQuery)
    check(err)
    part, err := mr.NextPart()
    for err == nil {
        if name := part.FormName(); name != "" {
            if part.FileName() != "" {
                fileInfos = append(fileInfos, cs.handleUpload(r, part))
            } else {
                r.Form[name] = append(r.Form[name], getFormValue(part))
            }
        }
        part, err = mr.NextPart()
    }
    return
}

func (cs *CsvServer) ListOfFiles(r render.Render){
    fi, err := ioutil.ReadDir(cs.Media)
    check(err)
    directory_list := []string{}
    for _, file := range fi{
        directory_list = append(directory_list, file.Name())
    }
    r.JSON(200, directory_list)
}

type Person struct {
    Name   string
    Age    int
    Emails []string
    Jobs   []*Job
}

type Job struct {
    Employer string
    Role     string
}

func (cs *CsvServer) Run(){
    m := martini.Classic()
    m.Use(martini.Static(cs.Static))

    m.Use(render.Renderer(render.Options{
        Directory: cs.Templates,
        Extensions: []string{".html"},
        Delims: render.Delims{"{[{", "}]}"},
    }))

    job1 := Job{Employer: "Monash", Role: "Honorary"}
    job2 := Job{Employer: "Box Hill", Role: "Head of HE"}

    person := Person{
        Name:   "jan",
        Age:    50,
        Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
        Jobs:   []*Job{&job1, &job2},
    }

    m.Get("/", func(r render.Render) {
        r.HTML(200, "index", person)
    })
    m.Get("/files", func(r render.Render) {
        cs.ListOfFiles(r)
    })
    m.Get("/upload/", func() string {
        return "hello world" // HTTP 200 : "hello world"
    })

    result := make(map[string][]*FileInfo, 1)
    m.Post("/upload/", func(w http.ResponseWriter, req *http.Request) {
        result["files"] = cs.handleUploads(req)
        b, err := json.Marshal(result)
        check(err)
        w.Header().Set("Cache-Control", "no-cache")
        jsonType := "application/json"
        if strings.Index(req.Header.Get("Accept"), jsonType) != -1 {
            w.Header().Set("Content-Type", jsonType)
        }

        fmt.Fprintln(w, string(b))
    })

    http.ListenAndServe(fmt.Sprintf(":%d", cs.Port), m)
    m.Run()
}
