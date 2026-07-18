package practice

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ============================================================
// io.Reader/Writer 练习（学完 review/io_reader.go 后做）
// ============================================================

// 练习 39：统计文件行数
// [TODO] CountLines 从 r 中读取并返回行数
// 提示: bufio.NewScanner(r).Scan()
func CountLines(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count, nil
}

// 练习 40：过滤并写入
// [TODO] FilterLines 从 r 读取每行，只保留包含 keyword 的行，写入 w
// 提示: bufio.Scanner + fmt.Fprintf(w, "%s\n", line)
func FilterLines(r io.Reader, w io.Writer, keyword string) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), keyword) {
			_, err := fmt.Fprintf(w, "%s\n", scanner.Text())
			if err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

// 练习 41：将多段数据拼接为一个 Reader
// [TODO] ConcatReaders 将多个 io.Reader 的内容拼接，返回新的 io.Reader
// 提示: io.MultiReader 一行就够
func ConcatReaders(readers ...io.Reader) io.Reader {
	return io.MultiReader(readers...)
}
