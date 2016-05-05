package goconf

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	char_end_of_line         = '\n'
	char_conf_prefix_section = '.'
	charset_conf_comments    = "#;"
	charset_conf_delimiter   = "=:"
)

// load configuration from file
// @param: filename, string file to load
func (cfg Configuration) loadFromFile(filename string) Configuration {
	fileStream, openErr := os.Open(filename)
	if nil != openErr {
		log.Fatal(fmt.Sprintf("Open file error, file=[%s], error=[%s]", filename, openErr.Error()))
		return cfg
	}
	defer fileStream.Close()

	buf := bufio.NewReader(fileStream)
	sectionStack := make([]Configuration, 10, 10)
	sectionDepth := 0
	sectionStack[sectionDepth] = cfg

	for {
		line, err := buf.ReadString(char_end_of_line)
		if io.EOF == err { // read end of the file
			break
		}
		if nil != err {
			message := fmt.Sprintf("Read file error, file-[%s], error=[%s]", filename, err.Error())
			panic(message)
			break
		}
		trimedLine := strings.Trim(line, " \n\t")
		charCount := len(trimedLine)

		if charCount < 1 { // Skip empty line
			continue
		}

		commentPos := strings.IndexAny(trimedLine, charset_conf_comments)
		if commentPos == 0 { // Skip comments
			continue
		}

		if charCount >= 3 && trimedLine[0] == '[' && trimedLine[charCount-1] == ']' { // process section
			sectionName := strings.TrimSpace(trimedLine[1 : charCount-1])
			secCharCount := len(sectionName)
			tempSecDepth := 1
			for pi := 0; pi < secCharCount; pi++ {
				if sectionName[pi] != char_conf_prefix_section {
					break
				}
				tempSecDepth++
			}
			if tempSecDepth == secCharCount {
				panic(fmt.Sprintf("Error section, file=[%s], section=[%s]", filename, sectionName))
				return cfg
			}
			if tempSecDepth > 1 {
				sectionName = string(sectionName[tempSecDepth-1:])
			}

			if tempSecDepth == sectionDepth+1 { // 下级配置
				newCfg := NewConfiguration()
				sectionStack[sectionDepth].Add(sectionName, newCfg)

				sectionDepth = tempSecDepth
				if len(sectionStack) > tempSecDepth {
					sectionStack[tempSecDepth] = newCfg
				} else {
					sectionStack = append(sectionStack, newCfg)
				}
				continue
			}

			if tempSecDepth <= sectionDepth { // 向下越级配置
				newCfg := NewConfiguration()
				sectionStack[tempSecDepth-1].Add(sectionName, newCfg)
				sectionStack[tempSecDepth] = newCfg
				sectionDepth = tempSecDepth
				continue
			}
			// 其他情况，越级错误
			panic(fmt.Sprintf("Parse Section error, file=[%s], section=[%s]", filename, sectionName))
			continue
		}

		delemitPos := strings.IndexAny(trimedLine, charset_conf_delimiter)
		if delemitPos > 0 { // process key-value
			optionName := strings.TrimSpace(trimedLine[0:delemitPos])
			optionValue := strings.TrimSpace(trimedLine[delemitPos+1:])
			sectionStack[sectionDepth][optionName] = optionValue
			continue
		}
	}
	return cfg
}

func (cfg Configuration) Read(filename string) Configuration {
	fileInfo, fileErr := os.Stat(filename)
	if nil != fileErr { // file not exist
		log.Fatal(fmt.Sprintf("Read file error, file=[%s], error=[%s]", filename, fileErr.Error()))
		return cfg
	}
	if fileInfo.IsDir() {
		return cfg
	} else {
		cfg.loadFromFile(filename)
	}
	return cfg
}
