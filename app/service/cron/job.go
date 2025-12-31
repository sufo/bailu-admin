/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package cron

import (
	"bailu/app/domain/entity"
	"bailu/app/domain/vo"
	"bailu/app/service/cron/jobs"
	"bailu/pkg/httpclient"
	"bailu/pkg/log"
	"bailu/utils"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/wire"
	"io"
	"net/http"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

var Inject2JobSet = wire.NewSet(wire.Struct(new(Inject2Jobs), "*"))

// 函数集合
var jobList map[string]FuncJob

// 暂时没有放开配置
const HttpExecTimeout = 300 //http任务执行时间

// job
type JobProxy interface {
	Run() (result string, err error)
}

// 本地job接口
type FuncJob interface {
	Name() string //返回函数名称描述
	Invoke(context.Context, map[string]any) (result string, err error)
}

// 函数任务
type ExecJob struct {
	ctx context.Context
	entity.Task
}

// 函数任务注入
type Inject2Jobs struct {
	Notice *jobs.ScheduleNotice
}

// 需要将定义的struct 添加到字典中；
// 字典 key 可以配置到 自动任务 调用目标 中；
func (j *Inject2Jobs) InitFuncJobs() {
	jobList = map[string]FuncJob{
		"ExamplesOne":    &jobs.ExamplesOne{},
		"ScheduleNotice": j.Notice,
		"CleanLog":       &jobs.CleanLogJob{},
	}
}

// 获取本地job函数列表
func GetFuncJobs() []*vo.KV {
	var options = make([]*vo.KV, 0)
	for k, v := range jobList {
		opt := &vo.KV{k, v.Name()}
		options = append(options, opt)
	}
	return options
}

// 函数任务执行
func (e *ExecJob) Run() (result string, err error) {
	var job = jobList[e.InvokeTarget]
	if job == nil {
		log.L.Warnf("[%s] ExecJob Run job nil", e.InvokeTarget)
		return
	}
	var params map[string]any = nil
	if e.Args != "" {
		params = make(map[string]any)
		err = json.Unmarshal([]byte(e.Args), &params)
		if err != nil {
			//参数错误，则不需要重试了，直接返回
			return "", err
		}
	}
	return job.Invoke(e.ctx, params)
}

// http任务
type HttpJob struct {
	entity.Task
}

// 处理了循环次数
//func (h *HttpJob) Run() (result string, err error) {
//	//默认执行一次
//	var count int8 = 0
//	var resp httpclient.ResponseWrapper
//LOOP:
//	if count < RetryTimes {
//		if h.task.HttpMethod == entity.TaskHTTPMethodGet {
//			//处理参数
//			var params map[string]any
//			err = json.Unmarshal([]byte(h.task.Args), &params)
//			if err != nil {
//				//参数错误，则不需要重试了，直接返回
//				return "", err
//			}
//			urlParams := utils.Map2String(params)
//			target := h.task.InvokeTarget
//			if strings.Contains(h.task.InvokeTarget, "?") {
//				target = "&" + urlParams
//			} else {
//				target = "?" + urlParams
//			}
//			resp = httpclient.Get(target, HttpExecTimeout)
//		} else {
//			//urlFields := strings.Split(task.InvokeTarget, "?")
//			//task.Command = urlFields[0]
//			//var params string
//			//if len(urlFields) >= 2 {
//			//	params = urlFields[1]
//			//}
//			resp = httpclient.PostParams(h.task.InvokeTarget, h.task.Args, HttpExecTimeout)
//		}
//		// 返回状态码非200，均为失败
//		if resp.StatusCode != http.StatusOK {
//			err = fmt.Errorf("HTTP状态码非200-->%d", resp.StatusCode)
//			count = count + 1
//			if count <= RetryTimes {
//				time.Sleep(time.Duration(count) * RetryInterval * time.Second)
//				goto LOOP
//			}
//		}
//	}
//
//	return resp.Body, err
//}

func (h *HttpJob) Run() (result string, err error) {
	var resp httpclient.ResponseWrapper
	target := h.InvokeTarget
	//if h.HttpMethod == entity.TaskHTTPMethodGet {
	if h.HttpMethod == "GET" {
		if h.Args != "" {
			//处理参数
			var params map[string]any
			err = json.Unmarshal([]byte(h.Args), &params)
			if err != nil {
				//这种情况应该是不用重试了
				return "", err
			}
			urlParams := utils.Map2UrlParams(params)
			if strings.Contains(h.InvokeTarget, "?") {
				target = "&" + urlParams
			} else {
				target = "?" + urlParams
			}
		}
		resp = httpclient.Get(target, HttpExecTimeout)
	} else {
		//urlFields := strings.Split(task.InvokeTarget, "?")
		//task.Command = urlFields[0]
		//var params string
		//if len(urlFields) >= 2 {
		//	params = urlFields[1]
		//}
		resp = httpclient.PostParams(h.InvokeTarget, h.Args, HttpExecTimeout)
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Body)
	}

	return resp.Body, err
}

// shell
type ShellJob struct {
	entity.Task
}

func (h *ShellJob) Run() (result string, err error) {
	//多个命令一行
	command := strings.Split(h.InvokeTarget, "\r\n")
	cmd := exec.Command("sh")
	// 定义一对输入输出流
	inReader, inWriter := io.Pipe()
	// 把输入流的给到命令行
	cmd.Stdin = inReader
	// 获取标准输入流和错误信息流
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	sizeIndex := len(command) - 1
	// 指定用户执行
	osUser, err := user.Lookup(command[sizeIndex])
	if err == nil {
		//log.Printf("uid=%s,gid=%s", osUser.Uid, osUser.Gid)
		uid, _ := strconv.Atoi(osUser.Uid)
		gid, _ := strconv.Atoi(osUser.Gid)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	}

	if err = cmd.Start(); err != nil {
		return "", err
	}
	// 正常日志
	go func() {
		logScan := bufio.NewScanner(stdout)
		for logScan.Scan() {
			text := logScan.Text()
			fmt.Println("std：", text)
		}
	}()
	// 错误日志
	go func() {
		scan := bufio.NewScanner(stderr)
		for scan.Scan() {
			s := scan.Text()
			log.L.Error("build error: ", s)
		}
	}()
	// 写指令
	go func() {
		lines := command[:sizeIndex]
		for i, str := range lines {
			_, err := inWriter.Write([]byte(str))
			if err != nil {
				fmt.Println(err)
			}
			_, err = inWriter.Write([]byte("\n"))
			if err != nil {
				fmt.Println(err)
			}
			if i == len(lines)-1 {
				_ = inWriter.Close()
			}
		}
	}()

	err = cmd.Wait()
	state := cmd.ProcessState
	//执行失败，返回错误信息
	if !state.Success() {
		return state.String(), err
	}
	return state.String(), err
}
