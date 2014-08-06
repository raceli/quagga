package main
 
import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
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
//	Age  int
}
 
func (b Location) String() string {
	return fmt.Sprintf("%s - %s", b.Protocol, b.Label)
}
 
func (a Add) String() string {
	return fmt.Sprintf("%s - %d", a.Network, a.Age )
}

func (r Remove) String() string {
	return fmt.Sprintf("%s", r.Network)
}
 
func main() {
	xmlFile, err := os.Open("Blackhole.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()
 
	b, _ := ioutil.ReadAll(xmlFile)
 
	var q Query
	xml.Unmarshal(b, &q)
 
	fmt.Println(q.Blackhole)

	for _, addNetwork := range q.AddList {
		fmt.Println("---")
		fmt.Printf("\t%s\n", addNetwork)
	}

	for _, removeNetwork := range q.RemoveList {
		fmt.Printf("\t%s\n", removeNetwork)
	}
}
