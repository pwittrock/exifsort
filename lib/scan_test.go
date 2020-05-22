package exifsort

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestSkipFileType(t *testing.T) {
	// Try just gobo.<suffix>
	for suffix := range mediaSuffixMap() {
		goodInput := fmt.Sprintf("gobo.%s", suffix)
		if skipFileType(goodInput) {
			t.Errorf("Expected True for %s\n", goodInput)
		}
	}
	// Try a simple upper case just gobo.<suffix>
	for suffix := range mediaSuffixMap() {
		goodInput := strings.ToUpper(fmt.Sprintf("gobo.%s", suffix))
		if skipFileType(goodInput) {
			t.Errorf("Expected True for %s\n", goodInput)
		}
	}

	// Try with many "." hey.gobo.<suffix>
	for suffix := range mediaSuffixMap() {
		goodInput := fmt.Sprintf("hey.gobo.%s", suffix)
		if skipFileType(goodInput) {
			t.Errorf("Expected True for %s\n", goodInput)
		}
	}

	badInput := "gobobob.."
	if skipFileType(badInput) == false {
		t.Errorf("Expected False for %s\n", badInput)
	}

	badInput = "gobo"
	if skipFileType(badInput) == false {
		t.Errorf("Expected False for %s\n", badInput)
	}

	// Try ".." at the end.<suffix>
	for suffix := range mediaSuffixMap() {
		badInput := fmt.Sprintf("gobo.%s..", suffix)
		if skipFileType(badInput) == false {
			t.Errorf("Expected False for %s\n", badInput)
		}
	}
}

func populateExifDir(t *testing.T, dir string, withExif bool, num int, fileno *int) {
	var readPath string

	if withExif {
		readPath = "../data/with_exif.jpg"
	} else {
		readPath = "../data/no_exif.jpg"
	}

	content, err := ioutil.ReadFile(readPath)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < num; i++ {
		targetPath := fmt.Sprintf("%s/file%d.jpg", dir, *fileno)

		*fileno++

		err := ioutil.WriteFile(targetPath, content, 0600)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testTmpDir(t *testing.T, parent string, name string) string {
	newDir, err := ioutil.TempDir(parent, name)
	if err != nil {
		t.Fatal(err)
	}

	return newDir
}

/*
	Root
	-with_exif
	  -nested_exif
	-no_exif
	-mixed_exif
*/
func buildTestDir(t *testing.T) string {
	fileNo := 0
	rootDir := testTmpDir(t, "", "root")
	exifDir := testTmpDir(t, rootDir, "with_exif")
	nestedDir := testTmpDir(t, exifDir, "nested_exif")
	noExifDir := testTmpDir(t, rootDir, "no_exif")
	mixedDir := testTmpDir(t, rootDir, "mixed_exif")

	populateExifDir(t, exifDir, true, 50, &fileNo)
	populateExifDir(t, noExifDir, false, 25, &fileNo)
	populateExifDir(t, mixedDir, true, 25, &fileNo)
	populateExifDir(t, mixedDir, false, 25, &fileNo)
	populateExifDir(t, nestedDir, true, 25, &fileNo)

	return rootDir
}

func TestScanDir(t *testing.T) {
	tmpPath := buildTestDir(t)
	defer os.RemoveAll(tmpPath)

	w := ScanDir(tmpPath, false)

	const correctNumInvalid uint64 = 50

	const correctNumValid uint64 = 100

	if correctNumInvalid != w.Invalid() {
		t.Errorf("Expected %d Invalid Count. Got %d\n",
			correctNumInvalid, w.Invalid())
	}

	if correctNumValid != w.Valid() {
		t.Errorf("Expected %d Valid Count. Got %d\n",
			correctNumValid, w.Valid())
	}
}
