package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

func main() {
	// 切分字符串之间的空格，获取每个单词
	input := "foo bar       \n     baz"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	fmt.Println()
	// 下面用例说明：
	// 1. split 函数返回 0， nil， nil 导致不断的从底层 reader 中获取数据
	// 2. 当 Scanner 中的初始化 buffer size（例子中是2） 不够之后，会自动扩容
	input = "abcdefghijgklmn"
	scanner = bufio.NewScanner(strings.NewReader(input))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		fmt.Printf("%t\t%d\t%s\n", atEOF, len(data), data)
		return 0, nil, nil
	}
	scanner.Split(split)
	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
	fmt.Println()
	// 通过 split 函数中的 error 中断 scan 操作
	input = "abcdefghijklmn"
	scanner = bufio.NewScanner(strings.NewReader(input))
	split = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		fmt.Printf("%t\t%d\t%s\n", atEOF, len(data), data)
		if atEOF {
			return 0, nil, errors.New("bad luck")
		}
		return 0, nil, nil
	}
	scanner.Split(split)
	buf = make([]byte, 12)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println("err: ", scanner.Err())
	}
	fmt.Println()
	//
	input = "foofoofoo"
	scanner = bufio.NewScanner(strings.NewReader(input))
	split = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if bytes.Equal(data[:3], []byte{'f', 'o', 'o'}) {
			return 3, []byte{'F'}, nil
		}
		if atEOF {
			return 0, nil, io.EOF
		}
		return 0, nil, nil
	}
	scanner.Split(split)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	fmt.Println()
	// ErrFinalToken 可以使得 Scan() 停止，并且不返回错误
	input = "foo end bar"
	scanner = bufio.NewScanner(strings.NewReader(input))
	split = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanWords(data, atEOF)
		if err == nil && token != nil && bytes.Equal(token, []byte{'e', 'n', 'd'}) {
			return 0, []byte{'E', 'N', 'D'}, bufio.ErrFinalToken
		}
		return
	}
	scanner.Split(split)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}
	fmt.Println()
	// 每次遍历返回的值是有个最大值的 => bufio.MaxScanTokenSize
	// 可以通过 Buffer() 设置 buffer 以及 token 的 size
	input = strings.Repeat("x", bufio.MaxScanTokenSize)
	scanner = bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println("err: ", scanner.Err())
	}
	//
	fmt.Println()
	input = "foo|bar"
	scanner = bufio.NewScanner(strings.NewReader(input))
	split = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		fmt.Println(atEOF, string(data))
		if i := bytes.IndexByte(data, '|'); i >= 0 {
			return i + 1, data[0:i], nil
		}
		if atEOF {
			return len(data), data[:len(data)], nil
		}
		return 0, nil, nil
	}
	scanner.Split(split)
	for scanner.Scan() {
		fmt.Println("content: ", scanner.Text())
	}
}
