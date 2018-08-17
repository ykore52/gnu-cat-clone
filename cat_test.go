package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func StubIO(inbuf string, fn func()) (string, string) {
	inr, inw, _ := os.Pipe()
	outr, outw, _ := os.Pipe()
	errr, errw, _ := os.Pipe()
	orgStdin := os.Stdin
	orgStdout := os.Stdout
	orgStderr := os.Stderr
	inw.Write([]byte(inbuf))
	inw.Close()
	os.Stdin = inr
	os.Stdout = outw
	os.Stderr = errw
	fn()
	os.Stdin = orgStdin
	os.Stdout = orgStdout
	os.Stderr = orgStderr
	outw.Close()
	outbuf, _ := ioutil.ReadAll(outr)
	errw.Close()
	errbuf, _ := ioutil.ReadAll(errr)

	return string(outbuf), string(errbuf)
}

func TestSimpleCat(t *testing.T) {
	_, stderr := StubIO("", func() {
		err := Cat("/etc/os-release", false, false, false, false, false, false, false)
		assert.Nil(t, err)
	})
	fmt.Println(stderr)
}

func TestCatFromStdin(t *testing.T) {
	const input_str = "Lorem ipsum dolor sit amet\n"
	stdout, _ := StubIO(input_str, func() {
		err := Cat("", true, false, false, false, false, false, false)
		assert.Nil(t, err)
	})
	assert.Equal(t, stdout, input_str)
}

func TestOptsNumber(t *testing.T) {

	// --------------------------------------------------------------------------
	if true {
		const input_str = "Lorem\nipsum\ndolor"
		const expected = "     1  Lorem\n     2  ipsum\n     3  dolor\n"
		stdout, _ := StubIO(input_str, func() {
			err := Cat("", true, true, false, false, false, false, false)
			assert.Nil(t, err)
		})
		assert.Equal(t, expected, stdout)
	}

	// --------------------------------------------------------------------------
	if true {
		const input_str = "Lorem\n\nipsum\ndolor"
		const expected = "     1  Lorem\n     2  \n     3  ipsum\n     4  dolor\n"
		stdout, _ := StubIO(input_str, func() {
			err := Cat("", true, true, false, false, false, false, false)
			assert.Nil(t, err)
		})
		assert.Equal(t, expected, stdout)
	}
}

func TestOptsNumberNonBlank(t *testing.T) {
	const input_str = "Lorem\n\nipsum\ndolor"
	const expected = "     1  Lorem\n\n     2  ipsum\n     3  dolor\n"
	stdout, _ := StubIO(input_str, func() {
		err := Cat("", true, true, true, false, false, false, false)
		assert.Nil(t, err)
	})
	assert.Equal(t, expected, stdout)
}

func TestOptsShowEnds(t *testing.T) {

	// --------------------------------------------------------------------------
	if true {
		const input_str = "Lorem\n\nipsum\ndolor"
		const expected = "Lorem$\n$\nipsum$\ndolor$\n"
		stdout, _ := StubIO(input_str, func() {
			err := Cat("", true, false, false, true, false, false, false)
			assert.Nil(t, err)
		})
		assert.Equal(t, expected, stdout)
	}

	// --------------------------------------------------------------------------
	if true {
		const input_str = "Lorem\n\nipsum\ndolor"
		const expected = "     1  Lorem$\n     2  $\n     3  ipsum$\n     4  dolor$\n"
		stdout, _ := StubIO(input_str, func() {
			err := Cat("", true, true, false, true, false, false, false)
			assert.Nil(t, err)
		})
		assert.Equal(t, expected, stdout)
	}

	// --------------------------------------------------------------------------
	if true {
		const input_str = "Lorem\n\nipsum\ndolor"
		const expected = "     1  Lorem$\n$\n     2  ipsum$\n     3  dolor$\n"
		stdout, _ := StubIO(input_str, func() {
			err := Cat("", true, true, true, true, false, false, false)
			assert.Nil(t, err)
		})
		assert.Equal(t, expected, stdout)
	}

}

func TestOptsShowNonprinting(t *testing.T) {

	if true {
		const input = "あいうえお"
		const expected = "M-cM-^AM-^BM-cM-^AM-^DM-cM-^AM-^FM-cM-^AM-^HM-cM-^AM-^J\n"

		stdout, _ := StubIO(input, func() {
			err := Cat("", true, false, false, false, true, false, false)
			assert.Nil(t, err)
		})
		assert.Equal(t, expected, stdout)
	}

	if true {
		const input = "あい\tうえお"
		const expected = "M-cM-^AM-^BM-cM-^AM-^D\tM-cM-^AM-^FM-cM-^AM-^HM-cM-^AM-^J\n"

		stdout, _ := StubIO(input, func() {
			err := Cat("", true, false, false, false, true, false, false)
			assert.Nil(t, err)
		})
		assert.Equal(t, expected, stdout)
	}

}

func TestOptsShowNonprintingAndShowTabs(t *testing.T) {

	const input = "あい\tうえお"
	const expected = "M-cM-^AM-^BM-cM-^AM-^D^IM-cM-^AM-^FM-cM-^AM-^HM-cM-^AM-^J\n"

	stdout, _ := StubIO(input, func() {
		err := Cat("", true, false, false, false, true, true, false)
		assert.Nil(t, err)
	})
	assert.Equal(t, expected, stdout)

}

func TestOptsShowTabs(t *testing.T) {

	const input = "あい\tうえお"
	const expected = "あい^Iうえお\n"

	stdout, _ := StubIO(input, func() {
		err := Cat("", true, false, false, false, false, true, false)
		assert.Nil(t, err)
	})
	assert.Equal(t, expected, stdout)

}

func TestOptsSqeezeBlank(t *testing.T) {

	const input = "A\nB\nC\n\nD\nE\n\n\nF\nG\n\n\n\nH\n\nI\n\n"
	const expected = "A\nB\nC\n\nD\nE\n\nF\nG\n\nH\n\nI\n\n"

	stdout, _ := StubIO(input, func() {
		err := Cat("", true, false, false, false, false, false, true)
		assert.Nil(t, err)
	})
	assert.Equal(t, expected, stdout)

}

func TestInputFromFile(t *testing.T) {

	fp, err := GetFileHandle("/etc/os-release")
	if err != nil {
		t.Fail()
	}
	fp.Close()

}

func TestCommandShortOptions(t *testing.T) {

	if true {
		args := []string{"-b"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.True(t, opt.is_number)
		assert.True(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"-e"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.True(t, opt.is_show_ends)
		assert.True(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"-n"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.True(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"-s"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.True(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"-t"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.True(t, opt.is_show_nonprinting)
		assert.True(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"-u"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"-v"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.True(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"-A"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.True(t, opt.is_show_ends)
		assert.True(t, opt.is_show_nonprinting)
		assert.True(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"-E"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.True(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"-T"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.True(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"-"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"/etc/os-release"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.False(t, opt.is_enable_stdin)
		assert.NotEmpty(t, opt.filename)
	}

	if true {
		args := []string{"-n", "/etc/os-release"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.True(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.False(t, opt.is_enable_stdin)
		assert.NotEmpty(t, opt.filename)
	}

}

func TestCommandLongOptions(t *testing.T) {
	if true {
		args := []string{"--show-all"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.True(t, opt.is_show_ends)
		assert.True(t, opt.is_show_nonprinting)
		assert.True(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"--number-nonblank"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.True(t, opt.is_number)
		assert.True(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"--show-ends"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.True(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"--number"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.True(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"--squeeze-blank"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.True(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"--show-tabs"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.False(t, opt.is_show_nonprinting)
		assert.True(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}

	if true {
		args := []string{"--show-nonprinting"}
		opt, return_code := read_options(args)
		assert.Zero(t, return_code)
		assert.False(t, opt.is_number)
		assert.False(t, opt.is_number_nonblank)
		assert.False(t, opt.is_show_ends)
		assert.True(t, opt.is_show_nonprinting)
		assert.False(t, opt.is_show_tabs)
		assert.False(t, opt.is_squeeze_blank)
		assert.True(t, opt.is_enable_stdin)
	}
}
