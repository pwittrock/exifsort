/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"flag"
	"fmt"
// For when we move to cobra
//	"github.com/matchstick/exifSort/cmd"
	"github.com/matchstick/exifSort/lib"
)

var filepathArg = ""

func main() {
	flag.StringVar(&filepathArg, "filepath", "", "File-path of image")
	flag.Parse()
	if filepathArg == "" {
		panic("Set filepath")
	}

	fmt.Println("Opening:", filepathArg)
	entry, err := exifSort.ExtractExifDate(filepathArg)
	if err != nil {
		panic(err)
	}
	if entry.Valid == false {
		fmt.Printf("No Exif Data\n")
		return
	}
	fmt.Printf("Retrieved %+v\n", entry)
//	cmd.Execute()
}
