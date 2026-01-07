/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package vo

import (
	"github.com/sufo/bailu-admin/global"
	"github.com/sufo/bailu-admin/utils"
	time2 "github.com/sufo/bailu-admin/utils/time"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"os"
	"runtime"
	"time"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

var ServerInfo = &Server{}

type Server struct {
	Runtime  Runtime  `json:"runtime"`  //运行时
	ServInfo ServInfo `json:"servInfo"` //系统信息
	Cpu      Cpu      `json:"cpu"`      //cpu
	Ram      Ram      `json:"ram"`      //内存
	Disk     Disk     `json:"disk"`     //磁盘
}

type Runtime struct {
	StartTime    string `json:"startTime"` //启动时间
	RunTime      string `json:"runTime"`   //运行时间
	GOOS         string `json:"goos"`      //目标操作系统
	NumCPU       int    `json:"numCpu"`    //当前cpu核数量
	Compiler     string `json:"compiler"`
	GoVersion    string `json:"goVersion"`
	NumGoroutine int    `json:"numGoroutine"` //当前存在的Go协程数
}

type ServInfo struct {
	Name string `json:"name"` //系统名称
	OS   string `json:"os"`   //操作系统
	Ip   string `json:"ip"`   //服务器ip
	Arch string `json:"arch"` //系统架构
}

type Cpu struct {
	Cpus  []float64 `json:"cpus"`
	Cores int       `json:"cores"`
}

type Ram struct {
	Used  string  `json:"used"`
	Total string  `json:"total"`
	Free  string  `json:"free"`
	Usage float64 `json:"usage"`
}

type Disk struct {
	Path   string  `json:"path"`
	FsType string  `json:"fsType"`
	Used   string  `json:"used"`
	Total  string  `json:"total"`
	Free   string  `json:"free"`
	Usage  float64 `json:"usage"`
}

func (s *Server) CopyTo(c *gin.Context) error {
	var err error = nil
	err = s.setCPU()
	if err != nil {
		return err
	}

	err = s.setServInfo(c)
	if err != nil {
		return err
	}
	err = s.setRAM()
	if err != nil {
		return err
	}
	err = s.setDisk()
	if err != nil {
		return err
	}

	s.setRuntime()
	return err
}

func (s *Server) setRuntime() {
	s.Runtime = Runtime{
		StartTime:    global.StartTime.Format(time2.CSTLayout),
		RunTime:      time2.FormatDuration(time.Since(global.StartTime)),
		GOOS:         runtime.GOOS,
		NumCPU:       runtime.NumCPU(),
		Compiler:     runtime.Compiler,
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
	}
}

func (s *Server) setCPU() error {
	cores, err := cpu.Counts(false)
	if err != nil {
		return err
	}
	cpus, err := cpu.Percent(time.Duration(200)*time.Millisecond, true)
	if err != nil {
		return err
	}
	//格式化保留小数
	for i, p := range cpus {
		cpus[i] = utils.Round(p, 1)
	}

	//fmt.Printf("CPU使用率: %.3f%% \n", cpus[0])
	s.Cpu = Cpu{cpus, cores}
	return nil
}

func (s *Server) setServInfo(c *gin.Context) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	s.ServInfo = ServInfo{
		hostname,
		runtime.GOOS,
		//utils.GetLocalIP(), //不准确
		c.ClientIP(),
		runtime.GOARCH,
	}
	return nil
}

func (s *Server) setRAM() error {
	u, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	s.Ram = Ram{
		fmt.Sprintf("%.2f", float64(u.Used)/MB),
		fmt.Sprintf("%.2f", float64(u.Total)/MB),
		fmt.Sprintf("%.2f", float64(u.Total-u.Used)/MB),
		utils.Round(u.UsedPercent, 1),
	}
	return nil
}

func (s *Server) setDisk() error {
	//u, err := disk.Usage("/")
	//获取文件路径
	//_, filename, _, _ := runtime.Caller(0)
	dir, _ := os.Getwd()
	u, err := disk.Usage(dir)
	if err != nil {
		return err
	}
	s.Disk = Disk{
		u.Path,
		u.Fstype,
		fmt.Sprintf("%.2f", float64(u.Used)/MB),
		fmt.Sprintf("%.2f", float64(u.Total)/MB),
		fmt.Sprintf("%.2f", float64(u.Total-u.Used)/MB),
		utils.Round(u.UsedPercent, 1),
	}
	return nil
}
