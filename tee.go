package main 

import(
	"io"
	"os"
	"log"
)

// TeeImpl MultiWriter with extra steps
func TeeImpl(in io.ReadSeeker, out []io.Writer) error {
	writer := io.MultiWriter(out...)

	buf := make([]byte, 256)
	if _, err := io.CopyBuffer(writer, in, buf); err != nil {
		return err
	}

	return nil
}

func main() {
	args := os.Args[1:]
	flags := os.O_CREATE | os.O_WRONLY
	perms := 0644

	var isAppend bool
	for _, element := range args {
		if element == "--append" {
			isAppend = true
		} 
	}
	
	if isAppend {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	
	var fds []io.Writer
	for _, element := range args {
		if element == "-" || element == "--append" {
			continue
		} else {					
			fd, err := os.OpenFile(element, flags, os.FileMode(perms))
			if err != nil {
				log.Fatal(err)
			}
			fds = append(fds, fd)	
		}
	}

	fds = append(fds, os.Stdout)
	TeeImpl(os.Stdin, fds)
}