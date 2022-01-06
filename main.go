package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"
	"golang.org/x/xerrors"
)

// "foo": "type-${foo}-${foo}"
// "number": ${hoge}
var pattern = regexp.MustCompile(`\${.*}`)

func main() {
	app := &cli.App{
		Name:   "jv",
		Usage:  "rich json validator",
		Action: Validate,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}

func Validate(c *cli.Context) error {
	if c.NArg() != 1 {
		cli.ShowAppHelp(c)
		return nil
	}
	file, err := os.Open(c.Args().Get(0))
	if err != nil {
		return xerrors.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	b := strings.Builder{}
	for scanner.Scan() {
		line := scanner.Text()
		// 0に置換することで文字列と数値に両方対応
		_, err := b.WriteString(pattern.ReplaceAllString(line, "0") + "\n")
		if err != nil {
			return xerrors.Errorf("failed to write string: %w", err)
		}
	}
	replaced := b.String()

	if err := json.Unmarshal([]byte(replaced), &map[string]interface{}{}); err != nil {
		var offset int64
		var e1 *json.UnmarshalTypeError
		var e2 *json.SyntaxError
		if xerrors.As(err, &e1) {
			offset = err.(*json.UnmarshalTypeError).Offset
		} else if xerrors.As(err, &e2) {
			offset = err.(*json.SyntaxError).Offset
		} else {
			fmt.Println(xerrors.Errorf("failed to unmarshal: %w", err))
			os.Exit(2)
		}
		failedLine := countNewLineOfBytes(replaced[0:offset]) + 1
		fmt.Printf("syntax error at line %d\n", failedLine)

		startOffset := offset - 10
		if startOffset < 0 {
			startOffset = 0
		}
		for i := startOffset; i < offset; i++ {
			if replaced[i] == '\n' {
				startOffset = i
			}
		}

		endOffset := offset + 10
		if endOffset > int64(len(replaced)) {
			endOffset = int64(len(replaced))
		}
		for i := offset; i < endOffset; i++ {
			if replaced[i] == '\n' {
				endOffset = i
				break
			}
		}

		fmt.Println(replaced[startOffset:endOffset])
		fmt.Println(strings.Repeat(" ", int(offset-startOffset)) + "^^^")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return nil
}

func countNewLineOfBytes(b string) int {
	count := 0
	for _, c := range b {
		if c == '\n' {
			count++
		}
	}
	return count
}
