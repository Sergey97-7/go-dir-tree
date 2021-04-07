package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func deepDirTree(output io.Writer, path string, printFiles bool, preString string) error {
	//open file/dir
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Fatal("err ", err)
	}

	//get info about file $f
	fi1, err1 := f.Stat()
	if err1 != nil {
		panic(err1)
	}

	if !fi1.IsDir() {
		if fi1.Size() == 0 {
			fmt.Printf("%v%v (%v)\n", "├───", fi1.Name(), "empty")
		} else {
			fmt.Printf("%v%v (%vb)\n", "└───", fi1.Name(), fi1.Size())
		}
	} else {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		dirs := make([]os.FileInfo, 0)
		for _, va := range files {
			if printFiles {
				dirs = append(dirs, va)
			} else if va.IsDir() {
				dirs = append(dirs, va)
			}
		}
		sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

		for i, v := range dirs {
			if v.IsDir() {
				var line = preString
				if i < len(dirs)-1 {
					fmt.Fprintf(output, "%v%v\n", preString+"├───", v.Name())
					line += "│\t"
				} else {
					fmt.Fprintf(output, "%v%v\n", preString+"└───", v.Name())
					line += "\t"
				}
				deepDirTree(output, path+string(os.PathSeparator)+v.Name(), printFiles, line)
			} else if printFiles {
				if i < len(dirs)-1 {
					if v.Size() == 0 {
						fmt.Fprintf(output, "%v%v (%v)\n", preString+"├───", v.Name(), "empty")
					} else {
						fmt.Fprintf(output, "%v%v (%vb)\n", preString+"├───", v.Name(), v.Size())
					}
				} else {
					if v.Size() == 0 {
						fmt.Fprintf(output, "%v%v (%v)\n", preString+"└───", v.Name(), "empty")
					} else {
						fmt.Fprintf(output, "%v%v (%vb)\n", preString+"└───", v.Name(), v.Size())
					}
				}
			}
		}
	}
	return err
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	err := deepDirTree(output, path, printFiles, "")
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
