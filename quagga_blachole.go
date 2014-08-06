package main

import (
    "net"
    "os"
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "strings"
    "bytes"
)

const (
    ReadBuffer int = 1024
)

type Query struct {
        Blackhole Location
        // Have to specify where to find episodes since this
        // doesn't match the xml tags of the data that needs to go into it
        AddList []Add `xml:"Add"`
        RemoveList []Remove `xml:"Remove"`
}

type Location struct {
        // Have to specify where to find the series title since
        // the field of this struct doesn't match the xml tag
        Protocol string
        Label    string
}

type Add struct {
        Network  string
        Age  int
}

type Remove struct {
        Network  string
//      Age  int
}

func (b Location) String() string {
        return fmt.Sprintf("%s - %s", b.Protocol, b.Label)
}

func (a Add) String() string {
        return fmt.Sprintf("network %s\n", a.Network )
}


func (r Remove) String() string {
        return fmt.Sprintf("no network %s\n", r.Network)
}

 
func main() {
    var configString string 
    configString = readConfig("/etc/quagga/bgpd.conf")
    fmt.Printf("%s\n", configString)

    strPassword := ""
    strEnablePassword := ""
    strEnable := "enable\n"
    strRead := "write terminal\n"
    strWrite := "write memory\n"
    strConfigure := "configure terminal\n"
    strQuit := "quit\n"
    strRouter := "router bgp 20093 view Blackhole\n"
//    addNetwork := "network 192.168.1.1/32\n"

    servAddr := "localhost:2605"

    tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
    if err != nil {
        println("ResolveTCPAddr failed:", err.Error())
        os.Exit(1)
    }
 
    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        println("Dial failed:", err.Error())
        os.Exit(1)
    }

    reply := make([]byte, 1024)
 
    _, err = conn.Read(reply)
    if err != nil {
        println("Write to server failed:", err.Error())
        os.Exit(1)
    }
 
    println("reply from server=", string(reply))

    var returnedStr string = writeReadIO(conn,strPassword)
    fmt.Printf("%s\n", returnedStr)
    var returnedCode int = checkForError(returnedStr)
    returnedStr = writeReadIO(conn,strEnable)
    returnedCode = checkForError(returnedStr)
    returnedStr = writeReadIO(conn,strEnablePassword)
    returnedCode = checkForError(returnedStr)
    returnedStr = writeReadIO(conn,strRead)
    returnedCode = checkForError(returnedStr)
    returnedStr = writeReadIO(conn,strConfigure)
    returnedCode = checkForError(returnedStr)
    returnedStr = writeReadIO(conn,strRouter)
    returnedCode = checkForError(returnedStr)
//    returnedStr = writeReadIO(conn,addNetwork)
//    println(string(returnedStr))
 

        xmlFile, err := os.Open("Blackhole.xml")
        if err != nil {
                fmt.Println("Error opening file:", err)
                return
        }
        defer xmlFile.Close()

        b, _ := ioutil.ReadAll(xmlFile)

        var q Query
        xml.Unmarshal(b, &q)

//      fmt.Println(q.Blackhole)

        for _, addNetwork := range q.AddList {
          var addString string = "network " + addNetwork.Network + "\n";
//          fmt.Printf("%s", x)
//          fmt.Printf("%s", addNetwork)
          returnedStr = writeReadIO(conn,addString)
          returnedCode = checkForError(returnedStr)
        }

        for _, removeNetwork := range q.RemoveList {
          var removeString string = "no network " + removeNetwork.Network + "\n";
//          fmt.Printf("%s", x)
//          fmt.Printf("%s", removeNetwork)
          returnedStr = writeReadIO(conn,removeString)
          returnedCode = checkForError(returnedStr)
//          println(string(returnedStr))
        }

    returnedStr = writeReadIO(conn,strRead)
    returnedCode = checkForError(returnedStr)
    returnedStr = writeReadIO(conn,strWrite)
    returnedCode = checkForError(returnedStr)
    returnedStr = writeReadIO(conn,strQuit)
    returnedCode = checkForError(returnedStr)
    returnedStr = writeReadIO(conn,strQuit)
    returnedCode = checkForError(returnedStr)
    fmt.Printf("%d\n", returnedCode)

    conn.Close()

}

func writeReadIO (conn net.Conn,  str string) string {

    
//    buf := bytes.NewBuffer(nil)

    var returnStr string = ""

    _, err := conn.Write([]byte(str))

    if err != nil {
        println("Write to server failed:", err.Error())
        os.Exit(1)
    }
    
    reply := make([]byte, ReadBuffer )
//    reply := make([]byte, 10 )
    _ , err = conn.Read(reply)

    if err != nil {
        println("Read from server failed:", err.Error())
        os.Exit(1)
    }
    returnStr = string(reply)

    return returnStr

}

func checkForError (str string) int {

   var returnVal int = 0
   b := []byte(str)
   n := bytes.Index(b, []byte{0})
   printStr := string(b[:n])
   strLen := len(printStr)
   fmt.Printf("*%s*\n",printStr)
   fmt.Printf("String Length %d\n",strLen)
   if ((strings.Index(printStr, "Password:") >= 0) && (strLen < 12))  { 
      fmt.Printf("Bad Password Entered - Exiting NOW !\n")
      returnVal = -1 
   } // Should never have to re-enter password
   if (strings.Index(printStr, "Unknown command:") >= 0) { 
      fmt.Printf("Unknown Command is and Indication that something is wrong - Exiting NOW !\n")
      returnVal = -1 
   } // Should never be in wrong place 

    if (returnVal < 0 ) {
       os.Exit(returnVal)
    }

   return returnVal

}
func readConfig (configFile string) string {
     configData , err  := ioutil.ReadFile(configFile)
    if err != nil {
        println("Read from Configuration File failed:", err.Error())
        os.Exit(1)
    }
     return string(configData)
}
