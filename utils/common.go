package utils

import (
	"math/rand"
	"syscall"
	"os/exec"
	"strings"
	"time"
	"io/ioutil"
	"fmt"
	"compress/gzip"
	"os"
	"bytes"
	"encoding/hex"
	"bufio"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}


func RunCommand(executable string, args ...string) (string, error) {
	cmd := exec.Command(executable, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytes)), nil
}

func CheckConfig(targetStr string) bool {
    file, err := os.Open("/proc/config.gz")
    if err != nil {
        return false
    }
    defer file.Close()

    gzReader, err := gzip.NewReader(file)
    if err != nil {
        return false
    }
    defer gzReader.Close()

    scanner := bufio.NewScanner(gzReader)
    target := []byte(targetStr)
    for scanner.Scan() {
        if bytes.Contains(scanner.Bytes(), target) {
			// fmt.Println(scanner.Text())
            return true
        }
    }
    return false
}

func FindBTFAssets() string {
    var utsname syscall.Utsname
    err := syscall.Uname(&utsname)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    btf_file := "a12-5.10-arm64_min.btf"
    if strings.Contains(B2S(utsname.Release[:]), "rockchip") {
        btf_file = "rock5b-5.10-arm64_min.btf"
    }
    fmt.Printf("Load btf_file=%s\n", btf_file)
    return btf_file
}

func B2S(bs []int8) string {
	ba := make([]byte, 0, len(bs))
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(bytes.TrimSpace(bytes.Trim(ba, "\x00")))
}

func HexStringToBytes(hexStr string) ([]byte, error) {
    // 去除可能的空格或前缀（如 "0x"）
    cleaned := strings.ReplaceAll(hexStr, " ", "")
    cleaned = strings.TrimPrefix(cleaned, "0x")

    // 解码为字节数组
    data, err := hex.DecodeString(cleaned)
    if err != nil {
        return nil, fmt.Errorf("invalid hex string: %v", err)
    }
    return data, nil
}