package review

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// ============================================================
// 补充模块：io.Reader/Writer 实战
// ============================================================
// 接口定义你已经知道（review/interfaces.go 里的 Reader/Writer），
// 这里演示标准库中常用的具体类型。

func BufferDemo() {
	fmt.Println("=== bytes.Buffer：内存读写 ===")

	// Buffer 同时实现 io.Reader 和 io.Writer
	var buf bytes.Buffer

	buf.WriteString("hello ")
	buf.WriteString("world")

	fmt.Printf("Buffer 内容: %q\n", buf.String()) // "hello world"
	fmt.Printf("可读字节: %d\n", buf.Len())

	// 读出来
	data := make([]byte, 5)
	n, _ := buf.Read(data)
	fmt.Printf("读取 %d 字节: %q\n", n, data[:n]) // "hello"
}

func ReaderDemo() {
	fmt.Println("\n=== strings.Reader / bytes.Reader：从已有数据生成 Reader ===")

	r := strings.NewReader("gopher1")
	b := make([]byte, 3)
	n, _ := r.Read(b)
	fmt.Printf("读取: %q\n", b[:n]) // "gop"

	n2, _ := r.Read(b)
	fmt.Printf("继续读: %q\n", b[:n2]) // "her"
}

func CopyDemo() {
	fmt.Println("\n=== io.Copy：流式拷贝 ===")

	src := strings.NewReader("学习 io.Copy —— 从头读到尾，自动分块")
	var dst bytes.Buffer

	n, _ := io.Copy(&dst, src)
	fmt.Printf("拷贝了 %d 字节: %q\n", n, dst.String())
}

func FmtIODemo() {
	fmt.Println("\n=== fmt.Fprint / fmt.Fscanf：格式化读写 ===")

	// Fprintf：写入 Writer
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "name=%s age=%d", "Alice", 30)
	fmt.Printf("Fprintf 写入: %q\n", buf.String())

	// Fscanf：从 Reader 中按格式读取
	r := strings.NewReader("Bob 25")
	var name string
	var age int
	fmt.Fscanf(r, "%s %d", &name, &age)
	fmt.Printf("Fscanf 读取: name=%q age=%d\n", name, age)
}

func ScannerDemo() {
	fmt.Println("\n=== bufio.Scanner：逐行扫描 ===")

	// 模拟一个"大文件"
	content := "第一行\n第二行\n第三行\n"
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("  line: %q\n", line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("扫描出错:", err)
	}
}

func PipeDemo() {
	fmt.Println("\n=== io.Pipe：Reader/Writer 管道（goroutine 间传递）===")

	// Pipe 创建一个同步的内存管道
	r, w := io.Pipe()

	// 写端在 goroutine 里写
	go func() {
		w.Write([]byte("hello from pipe"))
		w.Close() // 必须关闭，否则读端永远阻塞
	}()

	// 读端在主 goroutine 里读
	data, _ := io.ReadAll(r)
	fmt.Printf("管道传输: %s\n", data)
}
