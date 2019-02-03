package reader

import (
	"io/ioutil"
	"os"
)

/* Reading files requires checking most calls for erros.
   This helper will streamline our error checks below
*/
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Read file and return content as a string
func ReadFileContent(path string) string {
	data, err := ioutil.ReadFile(path)
	check(err)
	return string(data)
}

// Read file and return os.File object
func ReadFile(path string) *os.File {
	file, err := os.Open(path)
	check(err)
	return file
}
