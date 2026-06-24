//go:build !windows

package services

import "golang.org/x/sys/unix"

func diskInfo() DiskInfo {
	var stat unix.Statfs_t
	target := "/app"
	if err := unix.Statfs(target, &stat); err != nil {
		_ = unix.Statfs("/", &stat)
	}
	bsize := uint64(stat.Bsize)
	total := uint64(stat.Blocks) * bsize
	free := uint64(stat.Bavail) * bsize
	d := DiskInfo{FreeBytes: free, TotalBytes: total, Healthy: free >= minRecommendedFreeB}
	if total > 0 {
		d.UsedPct = int(100 * (total - free) / total)
	}
	return d
}
