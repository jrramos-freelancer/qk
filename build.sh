#!/usr/bin/env bash

# Make sure work files are untracked to prevent commiting potentially sensitive information
# about company tooling and infrastructure.
untrack_work_files() {
    local work_path="commands/work/work.go"
    local work_test_path="commands/work/work_test.go"

    if [[ $(git ls-files -v $work_path) == H* ]]; then
        untrack_command="git update-index --assume-unchanged $work_path"
        echo "commands/work/work.go should not be tracked by git. Running \`$untrack_command\` to untrack it." >&2
        $untrack_command
    fi

    if [[ $(git ls-files -v $work_test_path) == H* ]]; then
        untrack_command="git update-index --assume-unchanged $work_test_path"
        echo "commands/work/work_test.go should not be tracked by git. Running \`$untrack_command\` to untrack it." >&2
        $untrack_command
    fi
}

test() {
    if [[ ! -x "test.sh" ]]; then
        echo "test.sh is not executable; run: sudo chmod +x test.sh" >&2
        exit 1
    fi

    ./test.sh
}

build() {
    test
    if [[ $? -eq 0 ]]; then
        echo -e "\nAll tests passed. Proceeding to build..."
    else
        echo -e "\nTests failed. Aborting build..." >&2
        exit 1
    fi

    go build

    if [[ $? -eq 0 ]]; then
        echo "Build successful."
    else
        echo "Build failed. Please check the errors above." >&2
    fi
}

main() {
    untrack_work_files
    build
}

main