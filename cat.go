package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
)

//
// オプション用構造体
//
type Options struct {
	is_number           bool
	is_number_nonblank  bool
	is_show_ends        bool
	is_show_nonprinting bool
	is_show_tabs        bool
	is_squeeze_blank    bool
	is_enable_stdin     bool
	filename            string
}

//
// コンストラクタ
//
func NewOptions() *Options {
	s := new(Options)

	s.is_number = false
	s.is_number_nonblank = false
	s.is_show_ends = false
	s.is_show_nonprinting = false
	s.is_show_tabs = false
	s.is_squeeze_blank = false
	s.is_enable_stdin = false
	s.filename = ""

	return s
}

func (o Options) String() string {
	return fmt.Sprintf("%#v", o)
}

//*****************************************************************************

//
// オプションに従い、ファイルもしくは標準入力の内容を表示する
//
func Cat(
	filename string,
	is_enable_stdin bool,
	is_number bool,
	is_number_nonblank bool,
	is_show_ends bool,
	is_show_nonprinting bool,
	is_show_tabs bool,
	is_squeeze_blank bool) error {

	var fp *os.File

	if is_enable_stdin {
		fp = os.Stdin
	} else {
		fp, _ = GetFileHandle(filename)
		defer fp.Close()
	}

	reader := bufio.NewReader(fp)

	line_number := 0
	new_lines := 0

	for {

		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		var out_buf = make([]byte, 0, 1024)

		if is_number &&
			is_number_nonblank &&
			len(line) == 0 {
			// do nothing
		} else if is_number {
			line_number += 1
			out_buf = append(out_buf, []byte(fmt.Sprintf("%6d  ", line_number))...)
		}

		if len(line) != 0 {
			if is_show_nonprinting {
				for i := 0; i < len(line); i++ {
					char := &line[i]

					// 制御文字ではないとき
					if *char >= 32 {

						switch {
						// 通常文字は普通に出力
						case *char < 127:
							out_buf = append(out_buf, *char)

						// DEL は printable に置換
						case *char == 127:
							out_buf = append(out_buf, '^')
							out_buf = append(out_buf, '?')

						// ASCII 範囲外文字は別文字に置換
						default:
							out_buf = append(out_buf, 'M')
							out_buf = append(out_buf, '-')

							if *char >= 127+32 {

								if *char < 128+127 {
									out_buf = append(out_buf, *char-128)
								} else {
									out_buf = append(out_buf, '^')
									out_buf = append(out_buf, '?')
								}

							} else {

								out_buf = append(out_buf, '^')
								out_buf = append(out_buf, *char-128+64)

							}
						}

					} else if *char == '\t' && !is_show_tabs {
						out_buf = append(out_buf, '\t')
					} else if *char == '\n' {
						out_buf = append(out_buf, '\n')
					} else {
						out_buf = append(out_buf, '^')
						out_buf = append(out_buf, *char+64)
					}
				}

			} else { // if !is_show_nonprinting

				for i := 0; i < len(line); i++ {
					char := &line[i]
					if *char == '\t' && is_show_tabs {
						out_buf = append(out_buf, '^')
						out_buf = append(out_buf, *char+64)
					} else {
						out_buf = append(out_buf, *char)
					}
				}

			}

			new_lines = 0

		} else { // if len(line) == 0

			if is_squeeze_blank {
				new_lines += 1
				if new_lines >= 2 {
					continue
				}
			}

		}

		if is_show_ends {
			out_buf = append(out_buf, '$')
		}

		fmt.Println(string(out_buf))
	}

	return nil
}

//
// ファイルを開き、ファイルハンドルを返す
//
func GetFileHandle(name string) (*os.File, error) {

	var fp *os.File
	fp, err := os.Open(name)
	if err != nil {
		errors.Wrapf(err, "File could not be open %s: ", name)
		return nil, err
	}

	return fp, nil
}

//
// バージョン文字列の表示
//
func version() {
	fmt.Println("cat (GNU like implementation) 8.28")
}

//
// コマンド使用方法の表示
//
func usage() {
	fmt.Print(`Usage: cat [OPTION]... [FILE]...
  Concatenate FILE(s) to standard output.

  With no FILE, or when FILE is -, read standard input.

    -A, --show-all           equivalent to -vET
    -b, --number-nonblank    number nonempty output lines, overrides -n
    -e                       equivalent to -vE
    -E, --show-ends          display $ at end of each line
    -n, --number             number all output lines
    -s, --squeeze-blank      suppress repeated empty output lines
    -t                       equivalent to -vT
    -T, --show-tabs          display TAB characters as ^I
    -u                       (ignored)
    -v, --show-nonprinting   use ^ and M- notation, except for LFD and TAB
        --help     display this help and exit
        --version  output version information and exit

  Examples:
    cat f - g  Output f's contents, then standard input, then g's contents.
    cat        Copy standard input to standard output.
`)
}

//
// コマンドオプションの解釈
//
func read_options(osArgs []string) (*Options, int) {
	opt := NewOptions()
	for i := 0; i < len(osArgs); i++ {
		elem := osArgs[i]

		// 標準入力からの読み込みオプション
		if len(elem) == 1 && elem[0] == '-' {
			opt.is_enable_stdin = true
		} else if len(elem) >= 2 && elem[0] == '-' && elem[1] != '-' {

			// ショートオプションの解釈
			for _, rune := range elem[1:] {

				switch rune {
				case 'b':
					opt.is_number = true
					opt.is_number_nonblank = true
				case 'e':
					opt.is_show_ends = true
					opt.is_show_nonprinting = true
				case 'n':
					opt.is_number = true
				case 's':
					opt.is_squeeze_blank = true
				case 't':
					opt.is_show_tabs = true
					opt.is_show_nonprinting = true
				case 'u':
					/* We provide the -u feature unconditionally */
				case 'v':
					opt.is_show_nonprinting = true
				case 'A':
					opt.is_show_nonprinting = true
					opt.is_show_ends = true
					opt.is_show_tabs = true
				case 'E':
					opt.is_show_ends = true
				case 'T':
					opt.is_show_tabs = true

				default:
					fmt.Printf("cat: invalid option -- '%s'\n", elem)
					fmt.Printf("Try 'cat --help' for more information.\n")
					return nil, 1
				}

			}

		} else {

			// ファイル名, ロングオプションの解釈
			switch elem {
			case "--show-all":
				opt.is_show_nonprinting = true
				opt.is_show_ends = true
				opt.is_show_tabs = true
			case "--number-nonblank":
				opt.is_number = true
				opt.is_number_nonblank = true
			case "--show-ends":
				opt.is_show_ends = true
			case "--number":
				opt.is_number = true
			case "--squeeze-blank":
				opt.is_squeeze_blank = true
			case "--show-tabs":
				opt.is_show_tabs = true
			case "--show-nonprinting":
				opt.is_show_nonprinting = true
			case "--help":
				usage()
				return nil, 0
			case "--version":
				version()
				return nil, 0
			default:
				if byte(elem[0]) == '-' {
					fmt.Printf("cat: invalid option -- '%s'\n", elem)
					fmt.Printf("Try 'cat --help' for more information.\n")
					return nil, 1
				}
				opt.filename = elem
			}

		}
	}

	if opt.filename == "" && !opt.is_enable_stdin {
		opt.is_enable_stdin = true
	}

	return opt, 0
}

//
// main()
//
func main() {

	option, return_code := read_options(os.Args[1:])
	if return_code != 0 {
		os.Exit(return_code)
	}

	err := Cat(
		option.filename,
		option.is_enable_stdin,
		option.is_number,
		option.is_number_nonblank,
		option.is_show_ends,
		option.is_show_nonprinting,
		option.is_show_tabs,
		option.is_squeeze_blank)

	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
