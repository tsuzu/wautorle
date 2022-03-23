package runner

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Runner struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    *bufio.Reader
	stdoutRaw io.Closer
}

func New() (*Runner, error) {
	wordleCmd := "wordle"
	if cmd := os.Getenv("WORDLE_CMD"); cmd != "" {
		wordleCmd = cmd
	}

	cmd := exec.Command(wordleCmd)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Runner{
		cmd:       cmd,
		stdin:     stdin,
		stdout:    bufio.NewReader(stdout),
		stdoutRaw: stdout,
	}, nil
}

func (r *Runner) Close() error {
	r.stdin.Close()
	r.stdoutRaw.Close()
	r.cmd.Process.Kill()
	return r.cmd.Wait()
}

func (r *Runner) NextWord() (string, error) {
	result, err := r.stdout.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

type Color string

const (
	Green  Color = "green"
	Orange Color = "orange"
	Gray   Color = "gray"
)

type Char struct {
	Char  byte
	Color Color
}

type Result [5]Char

func (r Result) Finished() bool {
	for _, c := range r {
		if c.Color != Green {
			return false
		}
	}

	return true
}

func (r Result) String() string {
	buf := bytes.Buffer{}

	for _, c := range r {
		switch c.Color {
		case Green:
			buf.WriteString("G" + string(c.Char))
		case Orange:
			buf.WriteString("O" + string(c.Char))
		case Gray:
			buf.WriteString(string(c.Char))
		}
	}

	return buf.String()
}

func ParseResult(s string) Result {
	idx := 0
	var r Result
	var col Color = Gray
	for _, c := range []byte(s) {
		if c == 'G' {
			col = Green
			continue
		} else if c == 'O' {
			col = Orange
			continue
		}

		r[idx] = Char{
			Char:  c,
			Color: col,
		}
		col = Gray
		idx++
	}

	return r
}

func (r *Runner) WriteResult(result Result) error {
	_, err := r.stdin.Write([]byte(result.String() + "\n"))

	return err
}
