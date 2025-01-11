package config

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"os"
	"regexp"
)

const DefaultIgnoreFileListSize = 20
const DefaultIgnoreFileMapSize = 20

type IgnoreFile struct {
	Regex utils.StringBool `yaml:"regex"`
	File  string           `yaml:"file"`
}

type IgnoreFileCompile struct {
	Index      int
	IsRegex    bool
	StringFile string
	RegexFile  *regexp.Regexp
}

func (i *IgnoreFileCompile) CheckDirEntry(dir os.DirEntry) bool {
	return i.CheckName(dir.Name())
}

func (i *IgnoreFileCompile) CheckName(name string) bool {
	if i.IsRegex {
		return i.RegexFile.MatchString(name)
	} else {
		return name == i.StringFile
	}
}

func NewIgnoreFileCompile(ignore int, i *IgnoreFile) (*IgnoreFileCompile, error) {
	if i.Regex.IsEnable() {
		reg, err := regexp.Compile(i.File)
		if err != nil {
			return nil, err
		}

		return &IgnoreFileCompile{
			Index:      ignore,
			IsRegex:    true,
			StringFile: "",
			RegexFile:  reg,
		}, nil
	} else {
		return &IgnoreFileCompile{
			Index:      ignore,
			IsRegex:    false,
			StringFile: i.File,
			RegexFile:  nil,
		}, nil
	}
}

type IgnoreFileCompileList struct {
	Map map[int][]*IgnoreFileCompile
}

func (i *IgnoreFileCompileList) init() error {
	i.Map = make(map[int][]*IgnoreFileCompile, DefaultIgnoreFileMapSize)
	return nil
}

func (i *IgnoreFileCompileList) Add(ruleIgnore int, fileIgnore int, ifile *IgnoreFile) error {
	ignoreFile, err := NewIgnoreFileCompile(fileIgnore, ifile)
	if err != nil {
		return err
	}

	lst := i.Map[ruleIgnore]
	if lst == nil {
		lst = make([]*IgnoreFileCompile, 0, DefaultIgnoreFileListSize)
	}

	lst = append(lst, ignoreFile)
	i.Map[ruleIgnore] = lst
	return nil
}

type IgnoreForEachFunc func(ignoreFile *IgnoreFileCompile) (any, error)

func (i *IgnoreFileCompileList) ForEach(ruleIgnore int, fn IgnoreForEachFunc) (any, error) {
	lst := i.Map[ruleIgnore]
	if lst == nil {
		return nil, fmt.Errorf("rule not found")
	}

	for _, ignoreFile := range lst {
		res, err := fn(ignoreFile)
		if err != nil {
			return nil, err
		} else if res != nil {
			return res, nil
		}
	}
	return nil, nil
}
