package main

import (
	"fmt"
	"log"
	"os"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

func main() {
	fmt.Println("vim-go")
	fmt.Println("pid", os.Getpid())
	linuxproc.ReadMemInfo("/proc/meminfo")
	stat, err := linuxproc.ReadProcessStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail", err)
	}

	now := time.Now().Unix()
	elapse := now - int64(stat.Starttime)
	fmt.Println("elapse", elapse)
	// for _, s := range stat.CPUStats {
	// 	fmt.Println(s.User)
	// 	fmt.Println(s.Nice)
	// 	fmt.Println(s.System)
	// 	fmt.Println(s.Idle)
	// 	fmt.Println(s.IOWait)
	// }
}

// stat.CPUStatAll
// stat.CPUStats
// stat.Processes
// stat.BootTime
// ... etc
