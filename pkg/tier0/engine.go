package tier0

import (
	"strings"

	"safeexec/pkg/parser"
)

type Action string

const (
	ActionAllow Action = "ALLOW"
	ActionBlock Action = "BLOCK"
	ActionSkip  Action = "SKIP" 
)

type Rule interface {
	
	Evaluate(cmd *parser.ParsedCommand) (bool, Action, string)
}

type Engine struct {
	
	fastAllowlist map[string]bool
	
	commandRules map[string][]Rule
	
	
	globalRules []Rule
}

func NewEngine() *Engine {
	e := &Engine{
		fastAllowlist: make(map[string]bool),
		commandRules:  make(map[string][]Rule),
		globalRules:   make([]Rule, 0),
	}

	
	for _, cmd := range []string{"ls", "cd", "pwd", "whoami", "cat", "echo", "git"} {
		e.fastAllowlist[cmd] = true
	}

	
	e.commandRules["rm"] = []Rule{
		&ExactFlagPathRule{
			Flags:       []string{"-r", "-rf", "-fr"},
			TargetPaths: []string{"/", "/etc", "/var", "/*"},
			Message:     "This recursively deletes a critical system directory, which holds configurations and dependencies that cannot be undone. If you meant to remove one file, try `rm <specific-file>` instead.",
		},
	}

	e.commandRules["chmod"] = []Rule{
		&ChmodPermissiveRule{},
	}
	
	
	e.globalRules = append(e.globalRules, &NetworkToShellRule{})

	return e
}

func (e *Engine) Evaluate(cmd *parser.ParsedCommand) (Action, string) {
	
	if e.fastAllowlist[cmd.CommandName] {
		return ActionAllow, ""
	}

	
	if rules, exists := e.commandRules[cmd.CommandName]; exists {
		for _, rule := range rules {
			if matched, action, msg := rule.Evaluate(cmd); matched {
				return action, msg 
			}
		}
	}

	for _, rule := range e.globalRules {
		if matched, action, msg := rule.Evaluate(cmd); matched {
			return action, msg
		}
	}

	return ActionSkip, ""
}

type ExactFlagPathRule struct {
	Flags       []string
	TargetPaths []string
	Message     string
}

func (r *ExactFlagPathRule) Evaluate(cmd *parser.ParsedCommand) (bool, Action, string) {
	
	hasFlag := false
	for _, f := range r.Flags {
		if strings.Contains(cmd.Original, f) {
			hasFlag = true
			break
		}
	}
	if !hasFlag {
		return false, ActionSkip, ""
	}

	for _, p := range r.TargetPaths {
		
		if strings.HasSuffix(strings.TrimSpace(cmd.Original), p) {
			return true, ActionBlock, r.Message
		}
	}
	return false, ActionSkip, ""
}


type ChmodPermissiveRule struct{}

func (r *ChmodPermissiveRule) Evaluate(cmd *parser.ParsedCommand) (bool, Action, string) {
	if strings.Contains(cmd.Original, "777") || strings.Contains(cmd.Original, "a+rwx") {
		return true, ActionBlock, "Assigning overly permissive access rights (777/a+rwx) is forbidden."
	}
	return false, ActionSkip, ""
}


type NetworkToShellRule struct{}

func (r *NetworkToShellRule) Evaluate(cmd *parser.ParsedCommand) (bool, Action, string) {
	
	orig := strings.ToLower(cmd.Original)
	hasNetwork := strings.Contains(orig, "curl") || strings.Contains(orig, "wget")
	hasPipe := strings.Contains(orig, "|")
	hasShell := strings.Contains(orig, "bash") || strings.Contains(orig, "sh") || strings.Contains(orig, "zsh")

	if hasNetwork && hasPipe && hasShell {
		return true, ActionBlock, "Fetching a script from the network and piping it straight into a shell interpreter is prohibited. Download the file, inspect it, and execute it manually."
	}
	return false, ActionSkip, ""
}
