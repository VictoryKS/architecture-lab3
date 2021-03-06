package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strings"
    "io"
    "regexp"
)

func changeString(instr string) string{
  str := strings.SplitAfter(instr, ",")
  new := str[0]

  for i := 1; i < len(str); i++ {
    re := regexp.MustCompile(`^[a-zA-Z]`)
    re1 := regexp.MustCompile(`[a-zA-Z],$`)
    re2 := regexp.MustCompile(`^\s`)

    if ((re.FindString(str[i]) != "" || re1.FindString(str[i-1]) != "") && re2.FindString(str[i]) == "") {
      new += " " + str[i]
    } else {
      new += str[i]
    }
  }

  return new
}

func changeFile(indir, outdir, file_name string, c chan int){
  infile, err := os.Open(indir + "/" + file_name)
  if err != nil{
        log.Fatal(err)
        c <- -1
        os.Exit(1)
    }

      defer infile.Close()

      outfile_name := strings.Split(file_name, ".")[0] + ".res"

      if _, err := os.Stat(outdir); os.IsNotExist(err) {
        err = os.MkdirAll(outdir, 0755)
        if err != nil {
          log.Fatal(err)
          c <- -1
          os.Exit(1)
        }
      }

      outfile, err := os.Create(outdir + "/" + outfile_name)
      if err != nil{
            log.Fatal(err)
            c <- -1
            os.Exit(1)
        }

      defer outfile.Close()

    data := make([]byte, 64)
    var prevEnd byte
    var text string
    var prevLength int

    for {
        n, err := infile.Read(data)
        if err == io.EOF{
          text += string(data[prevLength - 1])
          outfile.WriteString(changeString(text))
          break
        }

        if len(text) != 0 {
          outfile.WriteString(changeString(text))
        }

        if prevEnd != 0 {
          text = string(prevEnd) + string(data[:n - 1])
        } else {
          text = string(data[:n - 1])
        }

        prevLength = n
        prevEnd = data[n - 1]
    }



  c <- 1
}

func main() {

  indir := os.Args[1]
  outdir := os.Args[2]

  files, err := ioutil.ReadDir(indir)
  if err != nil {
      log.Fatal(err)
  }

  c := make(chan int)

  for _, file := range files {
      if !file.IsDir() {
        go changeFile(indir, outdir, file.Name(), c)
      }
  }

  counter := 0
  len := len(files)

  for counter != len{
     if (<-c < 0) {
       len--
     } else {
       counter++
     }
  }

  fmt.Printf("Total number of processed files: %d.", counter)
}
