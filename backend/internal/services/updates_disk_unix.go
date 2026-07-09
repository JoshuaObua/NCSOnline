//go:build !windows

package services

import (
	"os"

	"golang.org/x/sys/unix"
)

func diskInfo() DiskInfo {
	var stat unix.Statfs_t
	target := firstNonEmpty(os.Getenv("APP_WORK_TREE"), "/app")
	if err := unix.Statfs(target, &stat); err != nil {
		_ = unix.Statfs("/", &stat)
	}
	bsize := uint64(stat.Bsize)
	total := uint64(stat.Blocks) * bsize
	free := uint64(stat.Bavail) * bsize
	d := DiskInfo{FreeBytes: free, TotalBytes: total, Healthy: free >= uint64Env("SMART_UPDATE_MIN_FREE_BYTES", 1<<30)}
	if total > 0 {
		d.UsedPct = int(100 * (total - free) / total)
	}
	return d
}
