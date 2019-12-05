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



func changeFile(indir, outdir, file_name string, c chan int){
  infile, err := os.Open(indir + "/" + file_name)
  if err != nil{
        log.Fatal(err)
        c <- -1
        os.Exit(1)
    }

      defer infile.Close()

    data := make([]byte, 64)
    text := ""

    for{
        n, err := infile.Read(data)
        if err == io.EOF{
            break
        }
      text = string(data[:n])
    }

    res := ""

  outfile_name := strings.Split(file_name, ".")[0] + ".res"
  //fmt.Println(outfile_name)

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

  outfile.WriteString(res)
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
