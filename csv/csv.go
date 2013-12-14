package csv

import (
    "bufio"
    "os"
    "log"
    "strings"
    "regexp"
    // "fmt"
)

const _EMAIL = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

var exp, _ = regexp.Compile(_EMAIL)

type CsvReader struct {
    Path string
    Output string
    ErrorOutput string
    Seperator string
    EmailColumn int
}

type Line []string

func (c *CsvReader) FindMailLike(line Line) (Line, bool){
    original := line
    found := false
    for key, val := range line{
        email := strings.ToLower(val)
        email_sep := strings.Split(email, " ")
        for _, ival := range email_sep{
            iemail := strings.ToLower(ival)
            if exp.MatchString(iemail){
                found = true
                temp := original[c.EmailColumn]
                original[c.EmailColumn] = iemail
                original[key] = temp
            }
        }
    }
    return original, found
}

func (c *CsvReader) Parse() (int, int, int, int){
    nr_of_supercools := 0
    found_and_replaced := 0
    nr_of_fucks := 0
    total := 0

    file, err := os.Open(c.Path)
    defer file.Close()
    if err != nil { panic(err) }

    fo, err := os.Create(c.Output)
    defer fo.Close()
    if err != nil { panic(err) }

    efo, err := os.Create(c.ErrorOutput)
    defer efo.Close()
    if err != nil { panic(err) }

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if line != ""{
            seperated := strings.Split(line, c.Seperator)
            email := strings.ToLower(seperated[c.EmailColumn])
            if exp.MatchString(email){
                nr_of_supercools++
                joined_line := strings.Join(seperated, ";")
                if _, err := fo.WriteString(joined_line); err != nil {
                    panic(err)
                }
                _, _ = fo.WriteString("\n")
            }else{
                nr_of_fucks++
                new_line, found := c.FindMailLike(seperated)
                if found{
                    found_and_replaced++
                    joined_line := strings.Join(new_line, ";")
                    if _, err := fo.WriteString(joined_line); err != nil {
                        panic(err)
                    }
                    _, _ = fo.WriteString("\n")
                }else{
                    joined_line := strings.Join(new_line, ";")
                    if _, err := efo.WriteString(joined_line); err != nil {
                        panic(err)
                    }
                    _, _ = efo.WriteString("\n")
                }
            }
        }
        total++
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return total, nr_of_supercools, found_and_replaced, nr_of_fucks
}
