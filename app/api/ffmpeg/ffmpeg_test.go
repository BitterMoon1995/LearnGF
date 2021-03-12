package ffmpeg

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestHello(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	//linux
	//cmd := exec.Command("sh")
	//windows
	command := exec.Command("cmd")
	command.Stdin = buffer

	go func() {
		buffer.WriteString("echo hello world > test.txt\n")
		buffer.WriteString("exit\n")
	}()

	if err := command.Run(); err != nil {
		fmt.Println(err)
		return
	}
}

func TestTransCode(t *testing.T) {
	cmdArguments := []string{"-rtsp_transport", "tcp",
		"-i", "rtsp://admin:RBYWEQ@192.168.3.67:554/h264/ch1/main/av_stream",
		"-codec:v", "libx264", "-map", "0", "-f", "hls", "-hls_list_size", "6", "-hls_wrap", "10",
		"-hls_time", "10", "D:/programs/nginx/nginx-1.18.0/html/hls/test.m3u8"}

	cmd := exec.Command("ffmpeg", cmdArguments...)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
	fmt.Printf("command output: %q", out.String())
}

func TestCreateDir(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("mkdir e:\\uav\\001\\2021\\3\\16\n")
	buffer.WriteString("exit\n")

	command := exec.Command("cmd")
	command.Stdin = buffer

	/*
		Output()就会执行命令，并返回执行结果([]byte)，
		后来发现返回的执行结果就是  Microsoft Windows [版本 10.0.18363.449]
		(c) 2019 Microsoft Corporation。保留所有权利。
		没有任何意义

		判断调用系统命令结果：
		命令执行出错时err != nil，err.Error()会报"exit status 1"
		原理：
		go语言中的log模块：
		Fatal函数会导致程序（调用os.Exit(1)）退出->退出返回值为1
		Panic函数会导致挂掉（且会打印出panic时的信息）并退出->退出返回值为2
	*/
	_, err := command.Output()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("命令执行出错")
	}
	//decoder := simplifiedchinese.GB18030.NewDecoder()
	//utf8Bytes, err := decoder.Bytes(output)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(utf8Bytes))
	//outputStr := string(output)

}

func TestConvertFile(t *testing.T) {

	cmdArgs := strings.Fields("-i" +
		" rtsp://admin:RBYWEQ@192.168.3.67:554/h264/ch1/main/av_stream" +
		" -vcodec copy -f mp4 E:/uav/mp4/nigger.mp4")
	command := exec.Command("ffmpeg", cmdArgs...)

	output, err := command.Output()
	if err != nil {
		fmt.Println("指令出错")
	} else {
		decoder := simplifiedchinese.GB18030.NewDecoder()
		utf8Bytes, _ := decoder.Bytes(output)
		log.Println(string(utf8Bytes))
	}
}

func fileDownload(request *ghttp.Request) {
	filename := "E:\\uav\\mp4\\nigger.mp4"

	file, _ := os.Open(filename)
	defer file.Close()

	fileHeader := make([]byte, 512)
	_, _ = file.Read(fileHeader)

	fileStat, _ := file.Stat()

	request.Response.Header().Set("Content-Disposition", "attachment; filename="+filename)
	request.Response.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	request.Response.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	_, _ = file.Seek(0, 0)
	buf := make([]byte, 1024*256) //256K缓冲区
	_, _ = io.CopyBuffer(request.Response.Writer, file, buf)

	return
}

func TestFileDownload(t *testing.T) {
	server := g.Server()

	server.BindHandler("/download", fileDownload)
	server.BindHandler("/str", func(request *ghttp.Request) {
		request.Response.Writeln("nigga~~~~~~~~~~~~~~~~~")
	})

	server.SetPort(9000)
	server.Run()
}
