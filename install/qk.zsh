: "${QK_ROOT:="${0:A:h:h}"}"
QK_BIN="${QK_BIN:-${QK_ROOT}/qk}"

qk() {
	local output
	output=$("${QK_BIN}" "$@") || {
		echo "No matching command found." >&2
		return 1
	}

	if [[ "$1" == "--debug" ]]; then
		echo "[ Debugging... ]"
		echo
		echo "$output" | sed 's/^/ > /'
		echo
	elif [[ -n "$output" ]]; then
		echo "[ Running... ]"
		echo
		echo "$output" | sed 's/^/ > /'
		echo
		eval "$output"
	fi
}
