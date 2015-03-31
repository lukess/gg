package ini

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ini is a map of maps in memory
type Conf map[string]map[string]string

var (
	conf = make(Conf)
)

// reading ini file line by line and maping to a map
func Init(fp string) *Conf {
	file, err := os.Open(fp)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//https://github.com/google/re2/wiki/Syntax
	section := regexp.MustCompile(`^\[(\w*)\]$`)
	kv := regexp.MustCompile(`^\s*(.*)\s*=\s*(.*)\s*`)
	var s string
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		// parse comments
		if strings.HasPrefix(t, "#") {
			continue
		}
		if strings.HasPrefix(t, ";") {
			continue
		}
		// parse section
		if section.MatchString(t) {
			matches := section.FindStringSubmatch(t)
			if len(matches) == 2 {
				s = matches[1]
				conf[s] = map[string]string{}
			} else {
				fmt.Println("section format error" + t)
			}
		}
		// parse key = value
		if s != "" {
			if kv.MatchString(t) {
				matches := kv.FindStringSubmatch(t)
				if len(matches) == 3 {
					conf[s][matches[1]] = matches[2]
				} else {
					fmt.Println("key=value format error:" + t)
				}
			}
		}
	}
	return &conf
}

func Dump() {
	for k, v := range conf {
		fmt.Printf("[%s]\n", k)
		for k1, v1 := range v {
			fmt.Printf("%s = %s\n", k1, v1)
		}
	}
}

func Set(params ...string) bool {
	if _, ok := conf[params[0]]; ok != true {
		conf[params[0]] = make(map[string]string)
	}
	switch len(params) {
	case 3:
		conf[params[0]][params[1]] = params[2]
		return true
	default:
		return false
	}
}

func Get(params ...string) string {
	if len(params) != 2 {
		fmt.Println("number of argumens for ini.Get() must be 2")
		return ""
	}
	return conf[params[0]][params[1]]
}

func GetAll(params ...string) []string {
	switch len(params) {
	case 0:
		var s []string
		for key := range conf {
			s = append(s, fmt.Sprintf("[%s]", key))
			for i, e := range conf[key] {
				s = append(s, fmt.Sprintf("%s=%s", i, e))
			}
		}
		return s
	case 1:
		var s []string
		for i, e := range conf[params[0]] {
			s = append(s, fmt.Sprintf("%s=%s", i, e))
		}
		return s
	case 2:
		return []string{conf[params[0]][params[1]]}
	default:
		return []string{""}
	}
}
