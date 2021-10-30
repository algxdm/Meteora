package lib

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	path   = strings.Split(os.Args[0], `\`)
	Path   = strings.Join(path[:len(path)-1], `\`)
	DlPath = fmt.Sprintf("%s\\dl", Path)
)

func catchError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func GetSelfIp() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	if err != nil {
		addrs, err := net.InterfaceAddrs()
		addr := strings.Split(addrs[1].String(), "/")[0]
		if err != nil {
			return "", nil
		}
		return addr, nil

	} else {
		addr, _, _ := net.SplitHostPort(conn.LocalAddr().String())
		return addr, nil
	}
}

// 发送文件
func SendFile(conn net.Conn, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if catchError(err) {
		return err
	}

	defer file.Close()

	// 获得文件信息
	fileInf, err := os.Stat(filePath)
	if catchError(err) {
		return err
	}

	// 发送文件信息(文件名, 文件大小)
	_, err = conn.Write([]byte(fmt.Sprintf(
		"%s\n%d",
		fileInf.Name(),
		fileInf.Size())))
	if catchError(err) {
		return err
	}

	// 校验是否收到信息
	chkData := make([]byte, 10)
	l, err := conn.Read(chkData)
	if string(chkData[:l]) != "suc" {
		return errors.New("Check failure ")
	}
	if catchError(err) {
		return err
	}

	data := make([]byte, 65536)
	for { // 循环发送数据
		dataLen, err := file.Read(data)
		data = data[:dataLen]
		if catchError(err) {
			return err
		} else {
			_, err = conn.Write(data)
			if catchError(err) {
				return err
			}
		}
		// 当读取数据大小小于
		if dataLen == 0 {
			break
		}
	}

	return nil
}

// 接收文件
func RecvFile(conn net.Conn) error {
	fInf := make([]byte, 100)

	l, err := conn.Read(fInf)
	if catchError(err) {
		return err
	}

	fInf = fInf[:l]
	fileInf := strings.Split(string(fInf), "\n")

	fileName := fileInf[0]
	fileSize, err := strconv.Atoi(fileInf[1])
	if catchError(err) {
		return err
	}

	file, err := os.OpenFile(
		fmt.Sprintf("%s\\"+fileName, DlPath),
		os.O_CREATE|os.O_WRONLY,
		0)
	if catchError(err) {
		return err
	}

	_, err = conn.Write([]byte("suc"))
	if catchError(err) {
		return err
	}

	defer file.Close()

	dataSize := 0

	data := make([]byte, 65536)
	for {
		recvDataSize, err := conn.Read(data)
		if catchError(err) {
			return err
		}
		data = data[:recvDataSize]

		_, err = file.Write(data)
		if catchError(err) {
			return err
		}

		dataSize += recvDataSize
		if dataSize >= fileSize {
			break
		}
	}

	return nil
}

// 自定义发送数据, 防止数据接收补全
func SendData(conn net.Conn, data []byte) error {
	dataLen := fmt.Sprintf("%d", len(data))
	_, err := conn.Write([]byte(dataLen))
	if catchError(err) {
		return err
	}

	chkData := make([]byte, 10)
	l, err := conn.Read(chkData)
	chkData = chkData[:l]
	if string(chkData) != dataLen {
		return errors.New("Send data failure ")
	}

	_, err = conn.Write(data)
	if catchError(err) {
		return err
	}
	return nil
}

// 接收数据
func RecvData(conn net.Conn) ([]byte, error) {
	d := make([]byte, 10)
	l, err := conn.Read(d)
	d = d[:l]

	dataLen, err := strconv.Atoi(string(d))
	if catchError(err) {
		return nil, err
	}

	data := make([]byte, dataLen)

	chkData := []byte(strconv.Itoa(dataLen))
	_, err = conn.Write(chkData)
	if catchError(err) {
		return nil, err
	}

	l, err = conn.Read(data)
	if dataLen != l {
		return nil, errors.New("Data recv failure ")
	}
	if catchError(err) {
		return nil, err
	}

	data = data[:l]
	return data, nil
}

// 局域网广播, 使控制端发现
func Broadcast() {
	rAddr := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 60720,
	}

	conn, err := net.DialUDP("udp", nil, &rAddr)
	if catchError(err) {
		return
	}

	for {
		_, err := conn.Write([]byte("00296966A71EBEAC"))
		if catchError(err) {
			return
		}
		time.Sleep(time.Second * 2)
	}
}

