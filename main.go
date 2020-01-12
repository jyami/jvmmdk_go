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
  defer file.Close()
  
  fmt.Printf("%#04x\n", readU4(file))
  fmt.Printf("%#04d\n", readU2(file))
  fmt.Printf("%#04d\n", readU2(file))
}

func readU4(f *os.File) uint32 {
  return binary.BigEndian.Uint32(readBytes(f, 4))
}

func readU2(f *os.File) uint16 {
  return binary.BigEndian.Uint16(readBytes(f, 2))
}

func readBytes(f *os.File, size int) []byte {
  b := make([]byte, size)
  _, err := f.Read(b)
  if err != nil { fatal(err.Error()) }
  return b
}