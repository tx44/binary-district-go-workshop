// go run ./cmd/kwaygen/main.go -n 10 -m 10 -limit 6

package main

import (
  "bufio"
  "flag"
  "fmt"
  "math/rand"
  "os"
  "strconv"
)

var (
  n = flag.Int("n", 5, "number of files")
  m = flag.Int("m", 10, "number of numbers per file")
  limit = flag.Int("limit", 50, "limit of parallel execution")
)

type Result struct {
  Err error
  File string
}

func main() {
  flag.Parse()

  ch := make(chan Result, *n)
  sem := make(chan struct{}, *limit)

  for i := 0; i < *n; i++ {
    sem <- struct{}{}
    go genFile(sem, ch, i, *m)
  }
  for i := 0; i < *n; i++ {
    if r := <- ch; r.Err != nil {
      fmt.Printf("create file error: %v\n", r.Err)
    } else {
      fmt.Printf("created file %q\n", r.File)
    }
  }
}

func genFile(sem chan struct{}, ch chan Result, i, m int) {
  defer func() { <- sem }()

  name := fmt.Sprintf("file.%d", i)
  file, err := os.Create(name)
  if err != nil {
    ch <- Result{
      Err: err,
    }
    return
  }
  defer file.Close()

  buf := bufio.NewWriter(file)

  for j := 0; j < m; j++ {
    x := rand.Intn(1000)
    err := writeRandNumber(x, buf)
    if err != nil {
      ch <- Result{
        Err: err,
      }
      return
    }
  }

  ch <- Result{
    Err: buf.Flush(),
    File: name,
  }
}

type StringWriter interface {
  WriteString(string) (int, error)
}

func writeRandNumber(x int, dest StringWriter) error {
  s := strconv.Itoa(x)
  _, err := dest.WriteString(s)
  if err == nil {
    _, err = dest.WriteString("\n")
  }
  return err
}