package commands

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"fmt"

	"github.com/ggiamarchi/http-check/common"
)

// DfValue represents an amount of storage (either used or available) for a specific mountpoint
type DfValue struct {
	Ratio    float64 `json:"ratio"`
	Kilobyte int     `json:"kilobyte"`
	Megabyte int     `json:"megabyte"`
	Gigabyte int     `json:"gigabyte"`
}

// DfEntry represents an mountpoint present in the `df` output
type DfEntry struct {
	Mountpoint string  `json:"mountpoint"`
	Filesystem string  `json:"filesystem"`
	Used       DfValue `json:"used"`
	Available  DfValue `json:"available"`
}

// Df execute the command `df -h` and give the output in an array of DfEntry
func Df() (*[]DfEntry, error) {

	cmd := exec.Command("df", "-k")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")

	entries := make([]DfEntry, 0, len(lines)-1)

	for _, v := range lines {

		if v == "" {
			continue
		}

		line := strings.Fields(v)

		if strings.HasPrefix(line[0], "/") {

			var used, available int

			used, err = strconv.Atoi(line[2])
			if err != nil {
				return nil, err
			}

			available, err = strconv.Atoi(line[3])
			if err != nil {
				return nil, err
			}

			entry := DfEntry{
				Mountpoint: line[len(line)-1],
				Filesystem: line[0],
				Used: DfValue{
					Kilobyte: used,
					Megabyte: used / 1024,
					Gigabyte: used / 1024 / 1024,
					Ratio:    float64(used) / float64(used+available),
				},
				Available: DfValue{
					Kilobyte: available,
					Megabyte: available / 1024,
					Gigabyte: available / 1024 / 1024,
					Ratio:    float64(available) / float64(used+available),
				},
			}

			entries = append(entries, entry)

		}

	}

	return &entries, nil
}

// DfCheck checks whether a condition is satisfied or not
// The value is in kilobyte
// The operator can be "greaterThan" or "lowerThan"
// The field can be "used" or "available"
func DfCheck(mountpoint string, field string, op string, value int) (bool, *common.Error) {

	entries, err := Df()

	if err != nil {
		return false, &common.Error{Msg: err.Error(), Code: "INTERNAL"}
	}

	var entry *DfEntry

	for i := 0; i < len(*entries); i++ {
		if (*entries)[i].Mountpoint == mountpoint {
			entry = &(*entries)[i]
			break
		}
	}

	if entry == nil {
		return false, &common.Error{Msg: fmt.Sprintf("Mountpoint %s not found", mountpoint), Code: "CLIENT"}
	}

	if field == "available" {
		if op == "greaterThan" {
			if entry.Available.Kilobyte > value {
				return true, nil
			}
			return false, nil
		}
		if op == "lowerThan" {
			if entry.Available.Kilobyte < value {
				return true, nil
			}
			return false, nil
		}
		return false, &common.Error{Msg: "Invalid operator. The operator can be 'greaterThan' or 'lowerThan'", Code: "CLIENT"}
	} else if field == "used" {
		if op == "greaterThan" {
			if entry.Used.Kilobyte > value {
				return true, nil
			}
			return false, nil
		}
		if op == "lowerThan" {
			if entry.Used.Kilobyte < value {
				return true, nil
			}
			return false, nil
		}
		return false, &common.Error{Msg: "Invalid operator. The operator can be 'greaterThan' or 'lowerThan'", Code: "CLIENT"}
	}

	return false, &common.Error{Msg: "Unkwnon error", Code: "INTERNAL"}
}
