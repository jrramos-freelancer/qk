if [[ -n "${BASH_SOURCE[0]:-}" ]]; then
	QK_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
else
	: "${QK_ROOT:="${HOME}/custom_scripts/qk"}"
fi
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
