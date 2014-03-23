package main




import(
	"bufio"
	//"io"
	"os"
	"fmt"
)

var Order_matrix [4][3]int

// "creates" file and writes to it
func Write_to_file(order_matrix [4][3]int){
	fo, err := os.Create("output.backup")
    if err != nil { panic(err) }
    // close fo on exit and check for its returned error
    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()
	buf := make([]byte,4)
	for i := 0; i < 4; i++{
		buf[i] = byte(order_matrix[i][2])
	}
	fmt.Println(buf)
	// write a chunk
    if _, err := fo.Write(buf[:4]); err != nil {
        panic(err)
    }
}

func Read_from_file(){
	fi, err := os.Open("output.backup")
    if err != nil { panic(err) }
    // close fi on exit and check for its returned error
    defer func() {
        if err := fi.Close(); err != nil {
            panic(err)
        }
    }()
    r:= bufio.NewReader(fi)
    buf :=make([]byte,4)
	for i := 0; i < 4; i++{
		b, _ :=r.ReadByte()
		buf[i] = b
	}
	fmt.Println(buf)


}

func main(){
	Order_matrix[3][2] = 5
	Write_to_file(Order_matrix)

	Read_from_file()
	/*
	// open input file
    fi, err := os.Open("README.md")
    if err != nil { panic(err) }
    // close fi on exit and check for its returned error
    defer func() {
        if err := fi.Close(); err != nil {
            panic(err)
        }
    }()

	fo, err := os.Create("output.txt")
    if err != nil { panic(err) }
    // close fo on exit and check for its returned error
    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()

    // make a buffer to keep chunks that are read
    buf := make([]byte, 1024)
    for {
        // read a chunk
        n, err := fi.Read(buf)
        if err != nil && err != io.EOF { panic(err) }
        if n == 0 { break }

        // write a chunk
        if _, err := fo.Write(buf[:n]); err != nil {
            panic(err)
        }
    }
    */
}