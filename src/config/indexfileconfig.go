package config

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"os"
	"regexp"
)

const DefaultIndexFileListSize = 20
const DefaultIndexFileMapSize = 20

type IndexFile struct {
	Regex utils.StringBool `yaml:"regex"`
	File  string           `yaml:"file"`
}

func (i *IndexFile) setDefault() {
	i.Regex.SetDefaultDisable()
}

func (i *IndexFile) check() ConfigError {
	if i.File == "" {
		return NewConfigError("file is empty")
	}

	return nil
}

type IndexFileCompile struct {
	Index      int
	IsRegex    bool
	StringFile string
	RegexFile  *regexp.Regexp
}

func (i *IndexFileCompile) CheckDirEntry(dir os.DirEntry) bool {
	return i.CheckName(dir.Name())
}

func (i *IndexFileCompile) CheckName(name string) bool {
	if i.IsRegex {
		return i.RegexFile.MatchString(name)
	} else {
		return name == i.StringFile
	}
}

func NewIndexFileCompile(index int, i *IndexFile) (*IndexFileCompile, error) {
	if i.Regex.IsEnable() {
		reg, err := regexp.Compile(i.File)
		if err != nil {
			return nil, err
		}

		return &IndexFileCompile{
			Index:      index,
			IsRegex:    true,
			StringFile: "",
			RegexFile:  reg,
		}, nil
	} else {
		return &IndexFileCompile{
			Index:      index,
			IsRegex:    false,
			StringFile: i.File,
			RegexFile:  nil,
		}, nil
	}
}

type IndexFileCompileList struct {
	Map map[int][]*IndexFileCompile
}

func (i *IndexFileCompileList) init() error {
	i.Map = make(map[int][]*IndexFileCompile, DefaultIndexFileMapSize)
	return nil
}

func (i *IndexFileCompileList) Add(ruleIndex int, fileIndex int, ifile *IndexFile) error {
	indexFile, err := NewIndexFileCompile(fileIndex, ifile)
	if err != nil {
		return err
	}

	lst, ok := i.Map[ruleIndex]
	if !ok || lst == nil {
		lst = make([]*IndexFileCompile, 0, DefaultIndexFileListSize)
	}

	lst = append(lst, indexFile)
	i.Map[ruleIndex] = lst
	return nil
}

type IndexForEachFunc func(indexFile *IndexFileCompile) (any, error)

func (i *IndexFileCompileList) ForEach(ruleIndex int, fn IndexForEachFunc) (any, error) {
	lst, ok := i.Map[ruleIndex]
	if !ok || lst == nil {
		return nil, fmt.Errorf("rule not found")
	}

	for _, indexFile := range lst {
		res, err := fn(indexFile)
		if err != nil {
			return nil, err
		} else if res != nil {
			return res, nil
		}
	}
	return nil, nil
}
