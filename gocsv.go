package main

import (
	"flag"
	"path/filepath"
	// "github.com/savruk/gocsv/csv"
	// "fmt"
    "github.com/savruk/gocsv/server"
	// "io/ioutil"
)

var (
	seperator    string
	filename     string
	email_column int
	run_server   bool
	port         int
)

func init() {
	flag.StringVar(&filename, "filename", "", "give filename to do things!?")
	flag.StringVar(&filename, "f", "", "give filename to do things!?")
	flag.StringVar(&seperator, "seperator", "", "seperator!?")
	flag.StringVar(&seperator, "s", "", "seperator!?")
	flag.IntVar(&email_column, "emailcolumn", -1, "give filename to do things!?")
	flag.IntVar(&email_column, "c", -1, "give filename to do things!?")
	flag.BoolVar(&run_server, "runserver", false, "run server!?")
	flag.BoolVar(&run_server, "r", false, "run server!?")
	flag.IntVar(&port, "port", 8888, "give port numbers")
	flag.IntVar(&port, "p", 8888, "give port number ")
}

func main() {
	flag.Parse()
	// filename := os.Args[1]
	// seperator := os.Args[2]
	// email_column, _ := strconv.Atoi(os.Args[3])

	// path, _ := filepath.Abs(filename)
	// outfile, _ := filepath.Abs(fmt.Sprintf("%s_cools.csv", strings.Split(filename, ".")[0]))
	// errorfile, _:= filepath.Abs(fmt.Sprintf("%s_errors.csv", strings.Split(filename, ".")[0]))
	// reader := &csv.CsvReader{
	// 	Path:        path,
	// 	Output:      outfile,
	// 	ErrorOutput:      errorfile,
	// 	Seperator:   seperator,
	// 	EmailColumn: email_column,
	// }

	// total, cool, far, nrf := reader.Parse()
	// fmt.Println("Total: ", total, "Cool: ", cool, "FoundAndReplaced: ", far,
	// 	"nr_of_fucks: ", nrf)
	if run_server {
		template_path, _ := filepath.Abs("./server/templates")
		static_path, _ := filepath.Abs("./server/static")
		media_path, _ := filepath.Abs("./server/media")
		cs := &server.CsvServer{
			Port: port,
			Static: static_path,
			Media: media_path,
			Templates: template_path,
		}
		cs.Run()
	} else {

	}
}
