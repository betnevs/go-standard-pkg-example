package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

// 实现一个简单的 writer，测试无 buffer 的 writer
type MyWriter int

func (w *MyWriter) Write(p []byte) (n int, err error) {
	fmt.Println("len:", len(p), ", p:", string(p))
	return len(p), nil
}

// 实现一个抛出错误的 writer
type MyWriterErr int

func (w *MyWriterErr) Write(p []byte) (n int, err error) {
	fmt.Printf("Write: %q\n", p)
	return 0, errors.New("bomb")
}

type Writer1 int

func (w *Writer1) Write(p []byte) (n int, err error) {
	fmt.Printf("Writer1: %q\n", p)
	return len(p), nil
}

type Writer2 int

func (w *Writer2) Write(p []byte) (n int, err error) {
	fmt.Printf("Writer2: %q\n", p)
	return len(p), nil
}

func main() {
	fmt.Println("My writer:")
	w := new(MyWriter)
	w.Write([]byte("a"))
	w.Write([]byte("b"))
	w.Write([]byte("c"))
	w.Write([]byte("d"))
	fmt.Println("Buffer writer:")
	// 对 MyWriter 进行一层包装，得到新的 buffer writer。只有当 buffer 达到设置的 buffer size
	// 或者显示调用 Flush()，统一输入到原 writer 中。
	bw := bufio.NewWriterSize(w, 3) // 此处设置 buffer size 为 3
	bw.Write([]byte{'b'})
	bw.Write([]byte{'c'})
	bw.Write([]byte{'d'})
	err := bw.Flush()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	// 当 buffer writer 调用底层的 writer 出现错误时，则会将错误返回给 buffer writer，下次再进行 Write() 操作时，会直接报错。
	fmt.Println("Buffer writer occur err:")
	w2 := new(MyWriterErr)
	bw2 := bufio.NewWriterSize(w2, 3)
	n, err := (bw2.Write([]byte{'a'}))
	fmt.Println("1, n:", n, "err:", err)
	n, err = bw2.Write([]byte{'b'})
	fmt.Println("2, n:", n, "err:", err)
	n, err = bw2.Write([]byte{'c'})
	fmt.Println("3, n:", n, "err:", err)
	n, err = bw2.Write([]byte{'d'})
	fmt.Println("4, n:", n, "err:", err)
	n, err = bw2.Write([]byte{'e'})
	fmt.Println("5, n:", n, "err:", err)
	err = bw2.Flush()
	fmt.Println(err)
	fmt.Println()
	// 获取 buffer writer 中缓存的字节数 n。
	fmt.Println("Buffer writer get buffer size:")
	w3 := new(MyWriter)
	bw3 := bufio.NewWriterSize(w3, 3)
	fmt.Println("buffer size:", bw3.Buffered())
	bw3.Write([]byte{'a'})
	fmt.Println("buffer size:", bw3.Buffered())
	fmt.Println()
	// 当写入的内容大于了 buffer writer 的 size，那么将直接写入到底层的 writer，直接跳过 buffer writer。
	fmt.Println("Buffer writer skip:")
	w4 := new(MyWriter)
	bw4 := bufio.NewWriterSize(w4, 4)
	bw4.Write([]byte("abcdefghijklmnopq"))
	fmt.Println()
	// 多个 writer 复用同一个 buffer writer，减少内存分配。
	fmt.Println("Buffer writer reset:")
	wr1 := new(Writer1)
	bwr1 := bufio.NewWriterSize(wr1, 3)
	bwr1.Write([]byte("abc"))
	bwr1.Write([]byte("efg"))
	bwr1.Flush() // 使用 Reset() 之前一定要记住 Flush() 之前 buffer 中的数据，防止丢失还在缓存中的数据。
	wr2 := new(Writer2)
	bwr1.Reset(wr2)
	bwr1.Write([]byte("xzv"))
	bwr1.Flush()
	fmt.Println()
	// 写入 byte, rune, string 三种常用类型。
	fmt.Println("Buffer writer write byte, rune, string:")
	w5 := new(MyWriter)
	bw5 := bufio.NewWriterSize(w5, 10)
	fmt.Println("buffer size:", bw5.Buffered())
	bw5.WriteByte('a')
	fmt.Println("buffer size:", bw5.Buffered())
	bw5.WriteRune('我')
	fmt.Println("buffer size:", bw5.Buffered())
	bw5.WriteString("abc好")
	fmt.Println("buffer size:", bw5.Buffered())
	fmt.Println()
	// 使用 ReaForm 不断从 Reader 中读取数据。
	fmt.Println("Buffer writer ReadForm:")
	w6 := new(MyWriter)
	bw6 := bufio.NewWriterSize(w6, 3)
	s := strings.NewReader("onetwothree")
	bw6.ReadFrom(s)
	bw6.Flush()
}
