// https://medium.com/@as27/a-simple-beginners-tutorial-to-io-writer-in-golang-2a13bfefea02
package main

import (
  "fmt"
  "encoding/json"
  "io"
  "bytes"
)

type Person struct {
 Id int
 Name string
 Age int
}

// type Writer interface {
//     Write(p []byte) (n int, err error)
// }

func (p *Person) Write(w io.Writer) {
 b, _ := json.Marshal(*p)
 w.Write(b)
}

func main() {
  p := Person{1, "Dimon", 25}
  var b bytes.Buffer
  // fmt.Printf("%T\n", b)
  p.Write(&b) // b - структура у которой есть метод Write удовлетворяет интерфейсу io.Writer
  fmt.Printf("%s", b.String())
}
