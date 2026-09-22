//go:build !windows

package application

import "syscall"

func GetDiskUsage(path string) (total, used, avail int64, err error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, 0, 0, err
	}
	total = int64(stat.Blocks) * int64(stat.Bsize)
	free := int64(stat.Bfree) * int64(stat.Bsize)
	avail = int64(stat.Bavail) * int64(stat.Bsize)
	used = total - free
	return total, used, avail, nil
}
