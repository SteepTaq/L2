package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

// Представляет состояние выполнения текущей программы
type runningState struct {
	cancel context.CancelFunc
	pgid   int
}

// Хранит текущее состояние терминала
var current struct {
	state runningState
}

// builtin команды в мапе для удобства
var builtinMap = map[string]func(args []string, r io.Reader, w io.Writer) error{
	"cd":   builtinCd,
	"pwd":  builtinPwd,
	"echo": builtinEcho,
	"kill": builtinKill,
	"ps":   builtinPs,
}

// Представляет один из этапов пайплайна, который может быть либо builtin функцией,
// либо внешней командой
type pipelineStage struct {
	// сегмент выполняемой команды
	argv []string
	// Флаг, указывающий, является ли этап пайплайна builtin функцией
	isBuiltin bool
	// Функция builtin, которую нужно выполнить
	builtin func(args []string, r io.Reader, w io.Writer) error
	// Аргументы для builtin функции
	args []string
	// Входной поток для этапа
	in io.Reader
	// Выходной поток для этапа
	out io.Writer
	// external command, которую нужно выполнить
	cmd *exec.Cmd
	// Канал для ожидания завершения этапа
	doneCh chan error
}

// Возвращает слайс этапов пайплайна из разобранной команды
func buildPipelineStages(segments [][]string) []*pipelineStage {
	stages := make([]*pipelineStage, 0, len(segments))
	for i, argv := range segments {
		s := &pipelineStage{argv: argv}
		// Если команда является builtin, то устанавливаем флаг и функцию
		if fn, ok := builtinMap[argv[0]]; ok {
			s.isBuiltin = true
			s.builtin = fn
			s.args = argv[1:]
		}

		// Если это первый этап, то устанавливаем Stdin
		if i == 0 {
			s.in = os.Stdin
		}

		// Если это последний этап, то устанавливаем Stdout
		if i == len(segments)-1 {
			s.out = os.Stdout
		}

		// Иначе создаем pipe и коннектим предыдущий этап с текущим
		if i > 0 {
			pr, pw := io.Pipe()
			stages[i-1].out = pw
			s.in = pr
		}

		stages = append(stages, s)
	}
	return stages
}

// Запускает builtin в горутинах или внешние команды пайплайна. Возвращает
// первый process group id для пайплайна
func startPipelineStages(ctx context.Context, stages []*pipelineStage) (int, error) {
	firstPGID := 0

	for _, s := range stages {
		// Если этап является builtin, то запускаем его в горутине
		if s.isBuiltin {
			s.doneCh = make(chan error, 1)
			go func(st *pipelineStage) {
				err := st.builtin(st.args, st.in, st.out)

				if pw, ok := st.out.(*io.PipeWriter); ok {
					_ = pw.Close()
				}

				st.doneCh <- err
				close(st.doneCh)
			}(s)

			continue
		}

		// Иначе запускаем внешнюю команду
		cmd := exec.CommandContext(ctx, s.argv[0], s.argv[1:]...)
		cmd.Stdin = s.in
		cmd.Stdout = s.out
		cmd.Stderr = os.Stderr

		if cmd.SysProcAttr == nil {
			cmd.SysProcAttr = &syscall.SysProcAttr{}
		}

		cmd.SysProcAttr.Setpgid = true
		if firstPGID != 0 {
			cmd.SysProcAttr.Pgid = firstPGID
		}

		if err := cmd.Start(); err != nil {
			if pw, ok := s.out.(*io.PipeWriter); ok {
				_ = pw.Close()
			}
			return 0, err
		}

		if firstPGID == 0 {
			firstPGID = cmd.Process.Pid
		}

		s.cmd = cmd
	}

	return firstPGID, nil
}

// Ожидает завершения всех этапов пайплайна и возвращает первую ошибку, если она есть.
func waitForPipelineStages(stages []*pipelineStage) error {
	var firstErr error

	for _, s := range stages {
		if s.cmd != nil {
			// Если этап является внешней командой, то ожидаем ее завершения
			if err := s.cmd.Wait(); err != nil && firstErr == nil {
				firstErr = err
			}

			if pw, ok := s.out.(*io.PipeWriter); ok {
				_ = pw.Close()
			}
		} else if s.doneCh != nil {
			// Иначе ожидаем завершения builtin функции
			if err := <-s.doneCh; err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		for range sigCh {
			// Если есть процесс группа, то отправляем сигнал SIGINT в нее
			if current.state.pgid != 0 {
				_ = syscall.Kill(-current.state.pgid, syscall.SIGINT)
			}

			if current.state.cancel != nil {
				current.state.cancel()
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	// Бесконечный цикл для чтения команд из Stdin
	for {
		wd, _ := os.Getwd()
		fmt.Fprintf(os.Stdout, "%s$ ", wd)
		line, err := reader.ReadString('\n')
		if err != nil {
			// Ctrl+D - выход из программы
			if err == io.EOF {
				fmt.Fprintln(os.Stdout)
				return
			}

			fmt.Fprintln(os.Stderr, "read error:", err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if err := evalLine(line); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func evalLine(line string) error {
	// Разбиваем строку на сегменты на поиск пайплайнов
	segTexts := strings.Split(line, "|")
	for _, p := range segTexts {
		t := strings.TrimSpace(p)
		if t == "" {
			fmt.Fprintln(os.Stderr, "syntax error, unexpected token")
			return nil
		}
	}
	if len(segTexts) == 0 {
		return nil
	}

	// Каждый сегмент разбиваем на поля
	var segments [][]string
	segments = make([][]string, 0, len(segTexts))
	for _, seg := range segTexts {
		args := strings.Fields(seg)
		if len(args) == 0 {
			continue
		}

		segments = append(segments, args)
	}
	if len(segments) == 0 {
		return nil
	}

	// Если сегмент один, то проверяем, является ли он builtin командой
	if len(segments) == 1 {
		argv := segments[0]
		if fn, ok := builtinMap[argv[0]]; ok {
			return fn(argv[1:], os.Stdin, os.Stdout)
		}
	}

	// Иначе запускаем пайплайн
	return runPipeline(segments)
}

func runPipeline(segments [][]string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stages := buildPipelineStages(segments)

	firstPGID, err := startPipelineStages(ctx, stages)
	if err != nil {
		cancel()
		return err
	}

	current.state = runningState{cancel: cancel, pgid: firstPGID}

	err = waitForPipelineStages(stages)
	current.state = runningState{}
	return err
}

func builtinCd(args []string, _ io.Reader, _ io.Writer) error {
	var dir string
	if len(args) == 0 {
		dir = os.Getenv("HOME")
		if dir == "" {
			dir = "/"
		}
	} else {
		dir = args[0]
	}

	return os.Chdir(dir)
}

func builtinPwd(_ []string, _ io.Reader, w io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, wd)
	return err
}

func builtinEcho(args []string, _ io.Reader, w io.Writer) error {
	_, err := fmt.Fprintln(w, strings.Join(args, " "))
	return err
}

func builtinKill(args []string, _ io.Reader, _ io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("kill: missing pid")
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("kill: invalid pid: %v", err)
	}

	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("kill: %v", err)
	}

	return nil
}

func builtinPs(_ []string, _ io.Reader, w io.Writer) error {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return err
	}

	pids := make([]int, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		pids = append(pids, pid)
	}

	sort.Ints(pids)

	bw := bufio.NewWriter(w)
	defer bw.Flush()

	for _, pid := range pids {
		pidStr := strconv.Itoa(pid)
		var cmd string

		if cmdlineBytes, err := os.ReadFile(filepath.Join("/proc", pidStr, "cmdline")); err == nil {
			cmd = strings.ReplaceAll(string(cmdlineBytes), "\x00", " ")
		}
		if strings.TrimSpace(cmd) == "" {
			if commBytes, err := os.ReadFile(filepath.Join("/proc", pidStr, "comm")); err == nil {
				cmd = strings.TrimSpace(string(commBytes))
			} else {
				continue
			}
		}

		if _, err := fmt.Fprintf(bw, "%d %s\n", pid, strings.TrimSpace(cmd)); err != nil {
			return err
		}
	}

	return nil
}