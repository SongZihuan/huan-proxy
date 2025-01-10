package utils

import (
	"os"
)

func Restart(newArgs ...string) (*os.Process, error) {
	args := make([]string, 0, len(os.Args))
	copy(args, os.Args)

	if len(newArgs) != 0 {
		args = append(args, newArgs...)
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	args0, err := os.Executable()
	if err != nil {
		return nil, err
	}

	attr := &os.ProcAttr{
		Dir:   wd,                                         // 新进程的工作目录
		Env:   os.Environ(),                               // 新进程的环境变量列表
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}, // 对应标准输入，标准输出和标准错误输出,若为nil,表示该进程启动时file是关闭的
	}

	p, err := os.StartProcess(args0, args[1:], attr)
	if err != nil {
		return nil, err
	}

	return p, nil
}
