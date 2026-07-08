package cmd

import "strings"

type ksh struct{}

// Korn shell instance
var Ksh Shell = ksh{}

const kshHook = `
_direnv_hook() {
  local previous_exit_status=$?;
  vars="$("{{.SelfPath}}" export ksh)";
  trap -- '' SIGINT;
  eval "$vars";
  trap - SIGINT;
  return $previous_exit_status;
};
if [[ ";${PROMPT_COMMAND[*]:-};" != *";_direnv_hook;"* ]]; then
  if [[ "$(declare -p PROMPT_COMMAND 2>&1)" == "declare -a"* ]]; then
    PROMPT_COMMAND=(_direnv_hook "${PROMPT_COMMAND[@]}")
  else
    PROMPT_COMMAND="_direnv_hook${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
  fi
fi
`

func (sh ksh) Hook() (string, error) {
	return kshHook, nil
}

func (sh ksh) Export(e ShellExport) (string, error) {
	var out strings.Builder
	for key, value := range e {
		if value == nil {
			out.WriteString(sh.unset(key))
		} else {
			out.WriteString(sh.export(key, *value))
		}
	}
	return out.String(), nil
}

func (sh ksh) Dump(env Env) (string, error) {
	var out strings.Builder
	for key, value := range env {
		out.WriteString(sh.export(key, value))
	}
	return out.String(), nil
}

func (sh ksh) export(key, value string) string {
	return "export " + sh.escape(key) + "=" + sh.escape(value) + ";"
}

func (sh ksh) unset(key string) string {
	return "unset " + sh.escape(key) + ";"
}

func (sh ksh) escape(str string) string {
	return KshEscape(str)
}

/*
 * Escaping
 */

// See also Quoting section in ksh(1).
func KshEscape(str string) string {
	if str == "" {
		return "''"
	}
	in := []byte(str)
	out := ""
	i := 0
	l := len(in)

	backslash := func(char byte) {
		out += string([]byte{BACKSLASH, char})
	}

	literal := func(char byte) {
		out += string([]byte{char})
	}

	for i < l {
		char := in[i]
		switch char {
		case BACKSLASH, '$', BACKTICK, '"':
			backslash(char)
		default:
			literal(char)
		}
		i++
	}

	out = "\"" + out + "\""

	return out
}
