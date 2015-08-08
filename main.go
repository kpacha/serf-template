package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	//"io/ioutil"
	//"text/template"
	"encoding/json"
)

// exit codes
const (
	OK           = iota
	SYNTAX_ERROR = iota
	CMD_FAILED   = iota
)

type Directive struct {
	template string
	result   string
	command  string
	tags     []string
}

type Member struct {
	Name     string
	Addr     string
	Port     int
	Tags     map[string]string
	Status   string
	Protocol map[string]int
}

type Serf_Output struct {
	Members []Member
}

func main() {
	args_len := len(os.Args)
	directives_len := args_len - 1
	directives := make([]Directive, directives_len)
	// for each args, parse into directives
	for i := 1; i < args_len; i = i + 1 {
		// split it into parts
		// 1st part: path to template file
		// 2nd part: path to result file
		// 3nd part: command to execute, optional
		// remaining parts: filter tags
		parts := strings.Split(os.Args[i], ":")
		parts_len := len(parts)
		// check number
		if parts_len < 2 {
			os.Exit(SYNTAX_ERROR)
		}
		for i := 0; i < parts_len; i = i + 1 {
			if len(parts[i]) == 0 {
				os.Exit(SYNTAX_ERROR)
			}
		}
		// register directive
		directives[i-1] = Directive{
			template: parts[0],
			result:   parts[1],
		}
		if parts_len > 2 {
			directives[i-1].command = parts[2]
		}
		if parts_len > 3 {
			directives[i-1].tags = make([]string, parts_len-3)
			for j := 0; j < parts_len-3; j = j + 1 {
				directives[i-1].tags[j] = parts[3+j]
			}
		}
	}

	// render template for each directives
	for i := 0; i < directives_len; i = i + 1 {
		cmd_args := []string{"members", "-format", "json"}
		for j := 0; j < len(directives[i].tags); j = j + 1 {
			cmd_args = append(cmd_args, "-tag")
			cmd_args = append(cmd_args, directives[i].tags[j])
		}
		cmd := exec.Command("serf", cmd_args...)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			os.Exit(CMD_FAILED)
		}
		var serf_output Serf_Output
		err = json.Unmarshal(out, &serf_output)
		if err != nil {
			fmt.Println(err)
			os.Exit(CMD_FAILED)
		}
		fmt.Println(&serf_output)
		//		tpl = template.ParseFiles(directives[i].template)
		//		tpl.Execute
	}

	os.Exit(OK)
}
