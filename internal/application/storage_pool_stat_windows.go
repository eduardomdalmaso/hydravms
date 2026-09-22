//go:build windows

package application

import (
	"golang.org/x/sys/windows"
)

func GetDiskUsage(path string) (total, used, avail int64, err error) {
	ptr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, 0, 0, err
	}
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64
	err = windows.GetDiskFreeSpaceEx(ptr, &freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes)
	if err != nil {
		return 0, 0, 0, err
	}
	total = int64(totalNumberOfBytes)
	avail = int64(freeBytesAvailable)
	used = total - int64(totalNumberOfFreeBytes)
	return total, used, avail, nil
}
