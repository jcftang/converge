#!/usr/bin/env bash
set -eo pipefail

ROOT=$(pwd)
TMP=$(mktemp -d -t converge.apply.XXXXXXXXXX)
function finish {
    rm -rf "$TMP"
}
trap finish EXIT

pushd "$TMP"

test_basic_conditionals() {
    SOURCE="$ROOT"/samples/conditional.hcl

    validate_unexecuted_branches() {
        if [ -f bar-output.txt ]; then
            echo "bar-output.txt exists but shouldn't"
            exit 1
        fi
        if [ -f baz-output.txt ]; then
            echo "baz-output.txt exists but shouldn't"
            exit 1
        fi
        return 0
    }

    "$ROOT"/converge apply --local "$SOURCE"

    validate_unexecuted_branches

    if [ -f foo-file.txt ]; then
        echo "foo-file.txt shouldn't exist!"
        exit 1
    fi

    if [ -f foo-file-2.txt ]; then
        echo "foo-file-2.txt shouldn't exist!"
        exit 1
    fi

    "$ROOT"/converge apply --local --only-show-changes -p "val=1" "$SOURCE"

    validate_unexecuted_branches

    if [ ! -f foo-file.txt ]; then
        echo "foo-file.txt doesn't exist!"
        exit 1
    fi

    if [ ! -f foo-file-2.txt ]; then
        echo "foo-file-2.txt doesn't exist!"
        exit 1
    fi

    if [[ "$(cat foo-file.txt)" != "foo1" ]]; then
        echo "foo-file.txt doesn't have the right content"
        echo "has '$(cat foo-file.txt)', want 'foo1'"
        exit 1
    fi

    if [[ "$(cat foo-file-2.txt)" != "foo-file.txt" ]]; then
        echo "foo-file-2.txt doesn't have the right content"
        echo "has '$(cat foo-file-2.txt)', want 'foo-file.txt'"
        exit 1
    fi
    return 0
}

test_language_conditionals() {
    SOURCE="$ROOT/samples/conditionalLanguages.hcl"

    test_lang_with_params() {
        params=${1:-""}
        expected=${2:-""}

        "$ROOT/converge" apply --local --only-show-changes -p "lang=${params}" "$SOURCE"

        if [ ! -f greeting.txt ]; then
            echo "greeting.txt doesn't exist!"
            exit 1
        fi

        if [[ "$(cat greeting.txt)" != "${expected}" ]]; then
            echo "greeting.txt doesn't have the right content"
            echo "has '$(cat greeting.txt)', want '${expected}'"
            exit 1
        fi
        return 0
    }

    test_lang_with_params "spanish" "hola"
    test_lang_with_params "french" "salut"
    test_lang_with_params "japanese" "もしもし"
    test_lang_with_params "english" "hello"
    test_lang_with_params "esperanto" "hello"

    return 0
}

test_basic_conditionals
test_language_conditionals

popd
