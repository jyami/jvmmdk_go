package main

import (
  "encoding/binary"
  "fmt"
  "os"
)

func main() {
  if len(os.Args) != 2 { fatal("not enough args.") } 
  readClassFile(os.Args[1])
}

func fatal(s string) {
  fmt.Println(s)
  os.Exit(1)
}

type ConstantPoolInfoTag byte

const (
  ConstantUTF8 ConstantPoolInfoTag = 1 + iota
  NotImplment_2
  NotImplment_3
  NotImplment_4
  NotImplment_5
  NotImplment_6
  ConstantClass
  ConstantString
  ConstantFieldRef
  ConstantMethodRef
  NotImplment_B
  ConstantNameAndType
  NotImplment_D
  NotImplment_E
)

var utf8Tbl map[int]string = map[int]string{}

func readClassFile(filename string) {
  file, err := os.Open(filename)
  if err != nil { fatal(err.Error()) }
  defer file.Close()
  
  magicNumber := readU4(file)
  fmt.Printf("magicNumber: %#04x\n", magicNumber)
  minorVersion := readU2(file)
  fmt.Printf("minorVersion: %#04d\n", minorVersion)
  majorVersion := readU2(file)
  fmt.Printf("majorVersion: %#04d\n", majorVersion)

  constantPoolCount := int(readU2(file))
  fmt.Printf("constantPoolCount: %d\n", constantPoolCount)

  for idx := 1; idx < constantPoolCount; idx++ {
    printConstantPoolRow(file, idx)
  }

  accessFlag := readU2(file)
  fmt.Printf("AccessFlag: %#04x\n", accessFlag)
  thisClassIdx := readU2(file)
  fmt.Printf("ThisClassIndex: %d\n", thisClassIdx)
  superClassIdx := readU2(file)
  fmt.Printf("SuperClassIndex: %d\n", superClassIdx)
  interfacesCount := int(readU2(file))
  fmt.Printf("InterfacesCount: %d\n", interfacesCount)

  for idx := 0; idx < interfacesCount; idx++ {
  }

  fieldsCount := int(readU2(file))
  fmt.Printf("FieldsCount: %d\n", fieldsCount)

  for idx := 0; idx < fieldsCount; idx++ {
  }

  methodsCount := int(readU2(file))
  fmt.Printf("MethodsCount: %d\n", methodsCount)

  for methodIdx := 0; methodIdx < methodsCount; methodIdx++ {
    fmt.Printf("method#%d ", methodIdx)
    printMethod(file)
    fmt.Printf("\n")
  }

  attributesCount := int(readU2(file))
  fmt.Printf("attributesCount: %d\n", attributesCount)

  for idx := 0; idx < attributesCount; idx++ {
    printAttribute(file)
  }
}

func printConstantPoolRow(file *os.File, idx int) {
  
  tag := ConstantPoolInfoTag(readU1(file))
  fmt.Printf("constantPoolRow#%d: ", idx)

  switch tag {
  case ConstantUTF8:

    length := int(readU2(file))
    s := string(readBytes(file, length))
    fmt.Printf("tag: ConstantUTF8[%#04x], length:%d content:%s\n", tag, length, s)
    utf8Tbl[idx] = s
  case ConstantClass:
    nameIndex := int(readU2(file))
    fmt.Printf("tag: ConstantClass[%#04x], nameIndex:%d\n", tag, nameIndex)
  case ConstantString:
    stringIndex := int(readU2(file))
    fmt.Printf("tag: ConstantString[%#04x], stringIndex:%d\n", tag, stringIndex)
  case ConstantFieldRef:
    classIndex := int(readU2(file))
    nameAndTypeIndex := int(readU2(file))
    fmt.Printf("tag: ConstantFieldRef[%#04x], classIndex:%d, nameAndTypeIndex:%d\n", tag, classIndex, nameAndTypeIndex)
  case ConstantMethodRef:
    classIndex := int(readU2(file))
    nameAndTypeIndex := int(readU2(file))
    fmt.Printf("tag: ConstantMethodRef[%#04x], classIndex:%d, nameAndTypeIndex:%d\n", tag, classIndex, nameAndTypeIndex)
  case ConstantNameAndType:
    nameIndex := int(readU2(file))
    descriptorIndex := int(readU2(file))
    fmt.Printf("tag: ConstantNameAndType[%#04x], nameIndex:%d, descriptorIndex:%d\n", tag, nameIndex, descriptorIndex)
  default:
    fatal(fmt.Sprintf("%d: not implemented %#04x", idx, tag))  
  }
}

func printMethod(file *os.File) {
  
  methodAccessFlag := readU2(file)
  fmt.Printf("methodAccessFlag: %#04x, ", methodAccessFlag)
  methodNameIndex := readU2(file)
  fmt.Printf("methodNameIndex: %d, ", methodNameIndex)
  descriptorIndex := readU2(file)
  fmt.Printf("descriptorIndex: %d, ", descriptorIndex)
  attributesCount := int(readU2(file))
  fmt.Printf("attributesCount: %d\n", attributesCount)

  for idx := 0; idx < attributesCount; idx++ {
    printAttribute(file)
  }
}

func printAttribute(file *os.File) {
  attributeNameIndex := int(readU2(file))
  fmt.Printf("  ")
  fmt.Printf("attributeNameIndex: %d, ", attributeNameIndex)
  attributeLength := int(readU4(file))
  fmt.Printf("attributeLength: %d, ", attributeLength)

  if utf8Tbl[attributeNameIndex] == "Code" {
    printCodeAttribute(file)
  } else {
    readBytes(file, attributeLength)
  }
}

type Operand struct {
  name string
  argc int
} 

var opTbl map[byte]Operand = map[byte]Operand{  
  0x12: Operand{ "ldc", 1 },
  0x2a: Operand{ "aload_0", 0 },
  0xb1: Operand{ "return", 0 },
  0xb2: Operand{ "getstatic", 2 },
  0xb6: Operand{ "invokevirtual", 2 },
  0xb7: Operand{ "invokespecial", 2 },
}

func printCodeAttribute(file *os.File) {
  maxStack := readU2(file)
  fmt.Printf("maxStack: %d, ", maxStack)
  maxLocals := readU2(file)
  fmt.Printf("maxLocals: %d, ", maxLocals)
  codeLength := int(readU4(file))
  fmt.Printf("codeLength: %d, ", codeLength)

  b := readBytes(file, codeLength)
  fmt.Printf("op: ")
  for idx := 0; idx < codeLength; idx++ {
    opCode := b[idx] 
    op := opTbl[opCode]
    fmt.Printf("< ")
    fmt.Printf("%s[%#04x] ", op.name, opCode)
    for argc := 0; argc < op.argc; argc++ {
      idx = idx + 1 
      fmt.Printf("%#04x ", b[idx])
    }
    fmt.Printf("> ")
  }
  fmt.Printf(", ")
  
  exceptionTableLength := int(readU2(file))
  fmt.Printf("exceptionTableLength: %d, ", exceptionTableLength)

  for idx := 0; idx < exceptionTableLength; idx++ {
    startPc := readU2(file)
    fmt.Printf("startPc: %d, ", startPc)
    endPc := readU2(file)
    fmt.Printf("endPc: %d, ", endPc)
    handlerPc := readU2(file)
    fmt.Printf("handlerPc: %d, ", handlerPc)
    catchType := readU2(file)
    fmt.Printf("catchType: %d, ", catchType)
  }

  attributesCount := int(readU2(file))
  fmt.Printf("attributesCount: %d, ", attributesCount)

  for idx := 0; idx < attributesCount; idx++ {
    printAttribute(file)
  }
}

func readU4(f *os.File) uint32 {
  return binary.BigEndian.Uint32(readBytes(f, 4))
}

func readU2(f *os.File) uint16 {
  return binary.BigEndian.Uint16(readBytes(f, 2))
}

func readU1(f *os.File) uint8 {
  return readBytes(f, 1)[0]
}

func readBytes(f *os.File, size int) []byte {
  b := make([]byte, size)
  _, err := f.Read(b)
  if err != nil { fatal(err.Error()) }
  return b
}