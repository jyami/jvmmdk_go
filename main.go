package main

import (
  "encoding/binary"
  "fmt"
  "os"
)

func main() {
  if len(os.Args) != 2 { fatal("not enough args.") } 
  readBinaryFile(os.Args[1])
}

func fatal(s string) {
  fmt.Println(s)
  os.Exit(1)
}

func readBinaryFile(filename string) {
  file, err := os.Open(filename)
  if err != nil { fatal(err.Error()) }
  
  fmt.Printf("%#04x\n", readU4(file))
  fmt.Printf("%#04d\n", readU2(file))
  fmt.Printf("%#04d\n", readU2(file))
}

func readU4(f *os.File) uint32 {
  b := make([]byte, 4)
  _, err := f.Read(b)
  if err != nil { fatal(err.Error()) }
  return binary.BigEndian.Uint32(b)
}

func readU2(f *os.File) uint16 {
  b := make([]byte, 2)
  _, err := f.Read(b)
  if err != nil { fatal(err.Error()) }
  return binary.BigEndian.Uint16(b)
}