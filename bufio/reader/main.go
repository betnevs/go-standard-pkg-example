package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type MyReader int

func (r *MyReader) Read(p []byte) (n int, err error) {
	fmt.Println("Do read")
	a := "abcdefghijklmnop"
	copy(p, a)
	return len(a), nil
}

func main() {
	// Buffer reader 会根据自己的size先从底层 reader 中读取数据然后缓存起来，
	// 如果真实读取的数量 size 大于 Buffer reader 中的 size，
	// 则会直接从底层 reader 中读取，相当于跳过 buffer reader。
	fmt.Println("Buffer reader do read:")
	rr := new(MyReader)
	bfrr := bufio.NewReader(rr)
	buf := make([]byte, 2)
	_, err := bfrr.Read(buf)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("%q\n", buf)
	}
	buf = make([]byte, 20)
	n, err := bfrr.Read(buf)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("read=%q, n=%d\n", buf[:n], n)
	}
	n, _ = bfrr.Read(buf)
	fmt.Printf("%q\n", buf[:n])
	fmt.Println()
	// Buffer reader 的 Peek() 方法帮助查看缓存中的前 n 个字节，而不是真的 advancing （可以理解成吃掉）。
	// 注意以下三种情况：
	// 1. 如果缓存不满，且数据少于 n 个字节，则尝试从 io.Reader 中读取
	// 2. 如果 Peek() 的数据量大于 buffer 的 size，将返回 bufio.ErrBufferFull
	// 3. 如果 Peek() 的数据量大于 buffer 的 size，将返回 EOF
	fmt.Println("Buffer reader seek:")
	s1 := strings.NewReader("abcdefghijklmnabcdefghijklmn")
	r := bufio.NewReaderSize(s1, 16)
	b, err := r.Peek(3)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("%q\n", b)
	}
	b, err = r.Peek(17)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("%q\n", b)
	}
	s2 := strings.NewReader("aaa")
	r.Reset(s2)
	b, err = r.Peek(10)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("%q\n", b)
	}
	fmt.Println()
	// 通过 Rest() 避免分配冗余的内存
	fmt.Println("Buffer reader reset:")
	s3 := strings.NewReader("abcde")
	r3 := bufio.NewReader(s3)
	by := make([]byte, 3)
	_, err = r3.Read(by)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("%q\n", by)
	}
	s4 := strings.NewReader("xyzavnaosdnvioasdnviansdoivnaiosdnvoiansdivnioasdv")
	r3.Reset(s4)
	_, err = r3.Read(by)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("%q\n", by)
	}
	_, err = r3.Read(by)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("%q\n", by)
	}
	fmt.Println()
	// Discard() 其实是跳过设置的字节数
	fmt.Println("Buffer reader discard:")
	r4 := new(MyReader)
	bfr := bufio.NewReader(r4)
	buf = make([]byte, 4)
	bfr.Read(buf)
	fmt.Printf("%q\n", buf)
	bfr.Discard(4)
	bfr.Read(buf)
	fmt.Printf("%q\n", buf)
	fmt.Println()
	// ReadByte() 读取一个字节，UnreadByte() 回退之前读取的字节
	fmt.Println("Buffer reader read/unread byte:")
	r5 := strings.NewReader("b我bcd")
	br5 := bufio.NewReader(r5)
	bt, err := br5.ReadByte()
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println(bt, br5.Buffered())
	}
	err = br5.UnreadByte()
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println(br5.Buffered())
	}
	bt, err = br5.ReadByte()
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println(bt, br5.Buffered())
	}
	fmt.Println()
	// ReadSlice() 存在几个需要注意的点：
	// 1. buffer reader 读取完底层的 reader 之后，并没有发现分割符号，则返回EOF
	// 2. buffer reader 的容量 size 小于底层 reader 的容量，则返回错误 ErrBufferFull
	fmt.Println("Buffer reader read slice:")
	s6 := strings.NewReader("abcdefghijklmnopqrstuvwxyz")
	bfr6 := bufio.NewReaderSize(s6, 16)
	token, err := bfr6.ReadSlice('|')
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("%q\n", token)
	}
	fmt.Println()
	// ReadBytes() 底层是使用了 ReadSlice()，不过它不受 buffer size 的限制，不会出现 ErrBufferFull
	fmt.Println("Buffer reader read bytes:")
	s7 := strings.NewReader(strings.Repeat("a", 20) + "|")
	bfr7 := bufio.NewReaderSize(s7, 16)
	token, err = bfr7.ReadBytes('|')
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("%q\n", token)
	}
	token, err = bfr7.ReadBytes('|')
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("%q\n", token)
	}
	fmt.Println()
	//
	fmt.Println("Buffer reader read line:")
	s8 := strings.NewReader(strings.Repeat("a", 20) + "\n" + "b\nc")
	bfr8 := bufio.NewReaderSize(s8, 16)
	token, isPrefix, err := bfr8.ReadLine()
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("Token: %q, prefix: %t\n", token, isPrefix)
	}
	token, isPrefix, err = bfr8.ReadLine()
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("Token: %q, prefix: %t\n", token, isPrefix)
	}
	token, isPrefix, err = bfr8.ReadLine()
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("Token: %q, prefix: %t\n", token, isPrefix)
	}
	token, isPrefix, err = bfr8.ReadLine()
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("Token: %q, prefix: %t\n", token, isPrefix)
	}
	fmt.Println()
	// 不断从 Reader 读入然后写入到 Writer 中。
	fmt.Println("Buffer reader write to:")
	bfr9 := bufio.NewReaderSize(new(R), 16)
	wn, err := bfr9.WriteTo(ioutil.Discard)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println("Written bytes:", wn)
	}
}

type R struct {
	n int
}

func (r *R) Read(p []byte) (n int, err error) {
	fmt.Printf("Read #%d\n", r.n)
	if r.n >= 10 {
		return 0, io.EOF
	}
	copy(p, "abcdefg")
	r.n += 1
	return 7, nil
}
