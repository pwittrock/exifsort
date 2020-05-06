package exifSort

import (
	"fmt"
	"github.com/dsoprea/go-exif-knife"
	"github.com/dsoprea/go-exif/v2"
	"github.com/dsoprea/go-exif/v2/common"
	"strconv"
	"strings"
	"time"
)

func formatError(label string, dateString string) (time.Time, error) {
	var t time.Time
	return t, fmt.Errorf("Bad Format for %s: %s Problem\n", dateString, label)
}

// Seconds are funny. The format may be "<sec> <milli>"
// or it may be with an extra decmial place such as <sec>.<hundredths>
func extractSecsFractionFromStr(secsStr string) (int, error) {
	splitSecs := strings.Split(secsStr, ".")
	if len(splitSecs) != 2 {
		return 0, fmt.Errorf("Not a fraction second")
	}

	// We only care about what is in front of the "."
	secs, err := strconv.Atoi(splitSecs[0])
	if err != nil {
		return 0, fmt.Errorf("Not a convertaible second")
	}
	return secs, nil
}

func extractTimeFromStr(exifDateTime string) (time.Time, error) {
	splitDateTime := strings.Split(exifDateTime, " ")
	if len(splitDateTime) != 2 {
		return formatError("Space Problem", exifDateTime)
	}
	date := splitDateTime[0]
	timeOfDay := splitDateTime[1]

	splitDate := strings.Split(date, ":")
	if len(splitDate) != 3 {
		return formatError("Date Split", exifDateTime)
	}

	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return formatError("Year", exifDateTime)
	}

	month, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return formatError("Month", exifDateTime)
	}

	day, err := strconv.Atoi(splitDate[2])
	if err != nil {
		return formatError("Day", exifDateTime)
	}

	splitTime := strings.Split(timeOfDay, ":")
	if len(splitTime) != 3 {
		return formatError("Time Split", exifDateTime)
	}

	hour, err := strconv.Atoi(splitTime[0])
	if err != nil {
		return formatError("Hour", exifDateTime)
	}

	minute, err := strconv.Atoi(splitTime[1])
	if err != nil {
		return formatError("Minute", exifDateTime)
	}

	second, err := strconv.Atoi(splitTime[2])
	if err != nil {
		second, err = extractSecsFractionFromStr(splitTime[2])
		if err != nil {
			return formatError("Sec", exifDateTime)
		}
	}

	t := time.Date(year, time.Month(month), day,
		hour, minute, second, 0, time.Local)
	return t, nil
}

type ExifDateEntry struct {
	Valid bool
	Path  string
	Time  time.Time
}

func ExtractExifDate(filepath string) (entry ExifDateEntry, err error) {
	var exifDateEntry ExifDateEntry
	exifDateEntry.Valid = false
	exifDateEntry.Path = filepath

	mc, err := exifknife.GetExif(filepath)
	if err != nil {
		return entry, err
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exifcommon.IfdPathStandard,
		"DateTime")
	if err != nil {
		return entry, err
	}

	_, found := mc.RootIfd.EntriesByTagId[it.Id]
	if found == false {
		return entry, err
	}

	ite := mc.RootIfd.EntriesByTagId[it.Id][0]
	value, err := ite.Value()
	if err != nil {
		return entry, err
	}

	exifDateEntry.Valid = true
	exifDateEntry.Time, err = extractTimeFromStr(value.(string))
	if err != nil {
		return exifDateEntry, err
	}

	return exifDateEntry, nil
}
