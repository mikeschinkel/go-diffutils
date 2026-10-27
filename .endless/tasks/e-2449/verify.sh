#!/usr/bin/env bash
# ── DO NOT EDIT ─────────────────────────────────────────────────────
# This suite belongs to E-2449 and records what was true when E-2449
# landed. Edit it only if you ARE E-2449. If your change breaks an
# assertion here, leave it alone — see .endless/tasks/CLAUDE.md.
source "$(dirname "${BASH_SOURCE[0]}")/../_harness.sh"

set -u

# Every module is checked as CI sees it: against its go.mod's requires, with
# no go.work.
export GOWORK=off

ROOT="$(git -C "${ENDLESS_VERIFY_DIR}" rev-parse --show-toplevel)" \
    || setup_error "cannot find the repository root"
cd "${ROOT}" || setup_error "cannot cd to ${ROOT}"

ROOTPATH="github.com/mikeschinkel/go-diffutils"
NESTED=(dmp sgdiff)
ALL_MODULES=(. dmp sgdiff)
LINTER="github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2"
BASE="$(git merge-base main HEAD)" || setup_error "cannot find the merge base with main"

# replaces prints one "old => new" line per replace directive in a go.mod
replaces() {
    go mod edit -json "$1/go.mod" | jq -r '.Replace // [] | .[] | "\(.Old.Path) => \(.New.Path)"'
}

# exported lists the exported top-level funcs and types of a module's
# non-test files at a commit, by name only
exported() {
    local rev=$1 dir=$2
    git ls-tree --name-only "${rev}" "${dir}/" | grep '\.go$' | grep -v '_test\.go$' \
        | while read -r f; do git show "${rev}:${f}"; done \
        | sed -nE 's/^func (\([^)]*\) )?([A-Z][A-Za-z0-9_]*).*/\2/p; s/^type ([A-Z][A-Za-z0-9_]*).*/\1/p' \
        | sort
}

section "Fail fast: every module vets and tests against its requires"
for dir in "${ALL_MODULES[@]}"; do
    if out=$(cd "${dir}" && go vet ./... 2>&1 && go test -race ./... 2>&1); then
        report_pass "${dir}: go vet and go test -race pass with GOWORK=off"
    else
        report_fail "${dir}: go vet and go test -race pass with GOWORK=off" "exit 0" "${out}"
    fi
done
if (( FAIL_COUNT > 0 )); then
    summary
fi

section "LICENSE is MIT plus go-difflib's full BSD 3-Clause notice (P1)"
license=$(cat LICENSE)
assert_contains "LICENSE opens with the MIT license" "MIT License" "$(head -n 1 LICENSE)"
assert_contains "LICENSE keeps Mike's copyright" "Copyright (c) 2026 Mike Schinkel" "${license}"
assert_contains "LICENSE names go-difflib's license as BSD 3-Clause" "BSD 3-Clause" "${license}"
assert_not_contains "LICENSE no longer says BSD 2-Clause" "2-Clause" "${license}"
assert_contains "LICENSE keeps Mezard's copyright line" "Copyright (c) 2013, Patrick Mezard" "${license}"
assert_contains "LICENSE has the source-retention condition" \
    "Redistributions of source code must retain the above copyright" "${license}"
assert_contains "LICENSE has the binary-reproduction condition" \
    "Redistributions in binary form must reproduce the above copyright" "${license}"
assert_contains "LICENSE has the non-endorsement condition" \
    "The names of its contributors may not be used to endorse or promote" "${license}"
assert_contains "LICENSE has the BSD disclaimer" "EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE" "${license}"

section "Source comments name go-difflib's license correctly (P2, P18, P20)"
header=$(sed -n '1,/^package /p' difflib.go)
assert_contains "difflib.go's header names BSD 3-Clause" "BSD 3-Clause" "${header}"
assert_not_contains "difflib.go's header no longer says MIT" "MIT" "${header}"
assert_contains "difflib.go's header keeps Mezard's copyright" "Copyright (c) 2013, Patrick Mezard" "${header}"
assert_contains "difflib.go's header points to LICENSE" "LICENSE file" "${header}"
assert_contains "difflib.go's header says the code was modified" "Modified by Mike Schinkel" "${header}"
pkgdoc=$(cat doc.go)
assert_contains "doc.go names go-difflib's license as BSD 3-Clause" "BSD 3-Clause" "${pkgdoc}"
assert_not_contains "doc.go no longer says MIT" "MIT" "${pkgdoc}"

section "README documents the modules, local development and licensing (P3)"
readme=$(cat README.md)
for dir in "${NESTED[@]}"; do
    assert_contains "README installs ${dir}" "go get ${ROOTPATH}/${dir}" "${readme}"
done
assert_contains "README installs root" "go get ${ROOTPATH}"$'\n' "${readme}"
assert_contains "README has a Local Development section" "## Local Development" "${readme}"
assert_contains "README tells you to run just work" "just work" "${readme}"
assert_contains "README explains GOWORK=off" "GOWORK=off" "${readme}"
assert_contains "README has a License section" "## License" "${readme}"
assert_contains "README credits go-difflib" "pmezard/go-difflib" "${readme}"

section "dmp and sgdiff require root by pseudo-version, with no replace (P4)"
for dir in "${NESTED[@]}"; do
    assert_eq "${dir}/go.mod has no replace directive" "" "$(replaces "${dir}")"
    required=$(go mod edit -json "${dir}/go.mod" \
        | jq -r --arg p "${ROOTPATH}" '.Require[] | select(.Path == $p) | .Version')
    if [[ "${required}" =~ ^v0\.0\.0-[0-9]{14}-([0-9a-f]{12})$ ]]; then
        sha=${BASH_REMATCH[1]}
        if git merge-base --is-ancestor "${sha}" main 2>/dev/null; then
            report_pass "${dir} requires root ${required}, a commit on main"
        else
            report_fail "${dir} requires root at a commit on main" "ancestor of main" "${sha}"
        fi
    else
        report_fail "${dir} requires root at a pseudo-version" "v0.0.0-<time>-<sha>" "${required:-none}"
    fi
    assert_contains "${dir}/go.sum records root ${required}" \
        "${ROOTPATH} ${required} h1:" "$(cat "${dir}/go.sum" 2>/dev/null)"
done

section "Every go.mod is tidy as CI sees it (P5)"
for dir in "${ALL_MODULES[@]}"; do
    if out=$(cd "${dir}" && go mod tidy -diff 2>&1); then
        report_pass "${dir}: go mod tidy -diff reports nothing"
    else
        report_fail "${dir}: go mod tidy -diff reports nothing" "no diff" "${out}"
    fi
done

section "Every module passes golangci-lint (P19)"
for dir in "${ALL_MODULES[@]}"; do
    if out=$(cd "${dir}" && go run "${LINTER}" run ./... --timeout=5m 2>&1); then
        report_pass "${dir}: golangci-lint reports nothing"
    else
        report_fail "${dir}: golangci-lint reports nothing" "0 issues" "${out}"
    fi
done
tests=$(go test -run 'ReportsFlushError' -v . 2>&1)
assert_contains "WriteUnifiedDiff returns its flush error" \
    "PASS: TestWriteUnifiedDiff_ReportsFlushError" "${tests}"
assert_contains "WriteContextDiff returns its flush error" \
    "PASS: TestWriteContextDiff_ReportsFlushError" "${tests}"

section "A justfile offers the recipes and finds its modules (P6)"
recipes=$(just --summary 2>&1)
for recipe in help work ensure-valid test lint fmt vet tidy build clean ci; do
    assert_contains "just offers ${recipe}" " ${recipe} " " ${recipes} "
done
assert_not_contains "the justfile does not list its modules" "dmp" "$(cat justfile)"

section "just work writes a go.work over every module, and git ignores it (P7, P21)"
scratch=$(mktemp -d) || setup_error "cannot make a scratch directory"
git archive HEAD | tar -x -C "${scratch}" || setup_error "cannot export HEAD"
# A module just work has never heard of must still land in go.work
mkdir -p "${scratch}/newmod" && printf 'module example.com/newmod\n\ngo 1.25.7\n' >"${scratch}/newmod/go.mod"
if out=$(cd "${scratch}" && GOWORK= just work 2>&1); then
    listed=$(cd "${scratch}" && GOWORK= go work edit -json | jq -r '.Use[].DiskPath' | sed 's|^\./||' | sort)
    expected=$(cd "${scratch}" && find . -name go.mod | sed 's|^\./||; s|/\{0,1\}go.mod$||; s|^$|.|' | sort)
    assert_eq "go.work uses every directory holding a go.mod" "${expected}" "${listed}"
else
    report_fail "just work succeeds" "exit 0" "${out}"
fi
rm -rf "${scratch}"
for f in go.work go.work.sum .endless/worktrees/x .endless/worktree.json; do
    assert_eq "git ignores ${f}" "${f}" "$(git check-ignore --no-index "${f}")"
done

section "test.yml runs without go.work and covers every module (P8)"
testyml=$(cat .github/workflows/test.yml)
assert_contains "test.yml sets GOWORK off" "GOWORK: 'off'" "${testyml}"
assert_contains "test.yml tests root" "go test -v -race -coverprofile" "${testyml}"
for dir in "${NESTED[@]}"; do
    assert_contains "test.yml tests ${dir}" "(cd ${dir} && go test" "${testyml}"
done
assert_contains "test.yml checks every go.mod is tidy" "go mod tidy -diff" "${testyml}"

section "release.yml makes directory-prefixed tags safely (P9–P11)"
rel=$(cat .github/workflows/release.yml)
modules_choice=$(sed -n '/^      module:/,/^      version:/p' .github/workflows/release.yml \
    | sed -n 's/^          - //p' | tr '\n' ' ')
assert_contains "module is a choice input" "type: choice" "${rel}"
assert_eq "module is a choice of root, dmp, sgdiff" "root dmp sgdiff " "${modules_choice}"
assert_contains "root tags vX.Y.Z" 'TAG="$VERSION"' "${rel}"
assert_contains "a nested module tags <dir>/vX.Y.Z" 'TAG="$MODULE/$VERSION"' "${rel}"
assert_contains "root's declared path is go-diffutils" "MODPATH=\"${ROOTPATH}\"" "${rel}"
assert_contains "release.yml checks the declared module path" '.Module.Path' "${rel}"
assert_contains "release.yml refuses a replace directive" 'contains replace directives' "${rel}"
assert_contains "release.yml refuses an existing tag" 'already exists' "${rel}"
assert_contains "release.yml tests with GOWORK off" "GOWORK: 'off'" "${rel}"
assert_contains "release.yml runs the module's tests" 'go test -v -race ./...' "${rel}"
assert_contains "a nested module's previous tag has its prefix" 'git tag -l "$MODULE/v*"' "${rel}"
assert_contains "root's previous tag is unprefixed" "git tag -l 'v*'" "${rel}"
assert_contains "nested notes list only commits touching its directory" 'PATHS=(-- "$DIR")' "${rel}"
assert_contains "root notes omit dmp, sgdiff and .endless" \
    "':(exclude)dmp' ':(exclude)sgdiff' ':(exclude).endless'" "${rel}"
assert_not_contains "release.yml names no go-dt module" "go-dt" "${rel}"

section "go-diffutils carries Endless's harness (P12)"
for f in _harness.sh _guard.sh; do
    assert_eq ".endless/tasks/${f} is committed" ".endless/tasks/${f}" \
        "$(git ls-files ".endless/tasks/${f}")"
done

section "Scope: doterr.go copies and exported identifiers stay (P15, P16)"
assert_eq "no doterr.go differs from main" "" \
    "$(git diff --name-only "${BASE}" HEAD -- doterr.go sgdiff/doterr.go)"
for dir in "${ALL_MODULES[@]}"; do
    assert_eq "${dir}: exported funcs and types are unchanged" \
        "$(exported "${BASE}" "${dir}")" "$(exported HEAD "${dir}")"
done

summary
