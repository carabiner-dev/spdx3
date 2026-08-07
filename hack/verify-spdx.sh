#!/usr/bin/env bash
# SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
# SPDX-License-Identifier: Apache-2.0
#
# Checks what this library writes against the SPDX project's own tooling.
#
# Every vendored document is read and written back out by the library, and a
# further document is built out of Go values rather than read from a file, so
# that the code turning Go values into SPDX ones is exercised too. Each result
# is put through two independent checks:
#
#   * the official SPDX 3.0.1 JSON schema, for structural conformance;
#   * the SPDX Java tools' Verify command, the reference implementation, for
#     semantic conformance.
#
# For the rendered documents both checks run over the input as well, and one is
# only reported as a failure when the input passes and our output does not.
# The authored document has no input to compare against and must simply be
# valid.
#
# Comparing against the input keeps a defect in an upstream example from
# failing this repository's builds: at the time of writing, spdx-examples'
# simplehtr document declares the same spdxId twice, and Verify rejects it
# whether they wrote it or we did.
#
#   ./hack/verify-spdx.sh [--keep DIR]
#
# Requires java and check-jsonschema on PATH. The schema and the tools jar are
# downloaded into a cache directory unless SPDX_JSON_SCHEMA and SPDX_TOOLS_JAR
# name local copies.

set -euo pipefail

readonly tools_version="2.0.7"
readonly tools_url="https://github.com/spdx/tools-java/releases/download/v${tools_version}/tools-java-${tools_version}.zip"
readonly schema_url="https://spdx.org/schema/3.0.1/spdx-json-schema.json"

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
readonly corpus="${root}/testdata/corpus"
cache="${SPDX_CONFORMANCE_CACHE:-${TMPDIR:-/tmp}/spdx3-conformance}"
rendered=""
keep=""

while [[ $# -gt 0 ]]; do
	case "$1" in
	--keep)
		keep="$2"
		shift 2
		;;
	-h | --help)
		sed -n '3,25p' "${BASH_SOURCE[0]}" | sed 's/^# \?//'
		exit 0
		;;
	*)
		echo "unknown argument: $1" >&2
		exit 2
		;;
	esac
done

for tool in java check-jsonschema; do
	if ! command -v "${tool}" >/dev/null; then
		echo "${tool} is required but not on PATH" >&2
		exit 1
	fi
done

mkdir -p "${cache}"

schema="${SPDX_JSON_SCHEMA:-${cache}/spdx-json-schema-3.0.1.json}"
if [[ ! -f "${schema}" ]]; then
	echo "Downloading the SPDX 3.0.1 JSON schema..."
	curl -sSL -o "${schema}" "${schema_url}"
fi

jar="${SPDX_TOOLS_JAR:-${cache}/tools-java-${tools_version}-jar-with-dependencies.jar}"
if [[ ! -f "${jar}" ]]; then
	echo "Downloading the SPDX Java tools ${tools_version}..."
	curl -sSL -o "${cache}/tools-java.zip" "${tools_url}"
	unzip -oq "${cache}/tools-java.zip" -d "${cache}"
fi

if [[ -n "${keep}" ]]; then
	rendered="${keep}"
	mkdir -p "${rendered}"
else
	rendered="$(mktemp -d)"
	trap 'rm -rf "${rendered}"' EXIT
fi

echo "Rendering the vendored documents through the library..."
(cd "${root}" && go run ./hack/render -authored -out "${rendered}" >/dev/null)

# schema_ok and verify_ok answer whether a document passes each check, and
# schema_output and verify_output say why it did not. None of them pipe the
# checker's output anywhere: a reader closing the pipe early would send it
# SIGPIPE, which under pipefail reads as the document having failed.
schema_output() {
	check-jsonschema --schemafile "${schema}" "$1" 2>&1 || true
}

schema_ok() {
	check-jsonschema --schemafile "${schema}" "$1" >/dev/null 2>&1
}

verify_output() {
	java -jar "${jar}" Verify "$1" 2>&1 || true
}

verify_ok() {
	local out
	out="$(verify_output "$1")"
	[[ "${out}" == *"This SPDX Document is valid."* ]]
}

# report prints the first lines of a checker's output, indented.
report() {
	local out="$1" lines=()
	mapfile -t lines <<<"${out}"
	printf '    %s\n' "${lines[@]:0:20}"
}

# run_check reports on one check of one document. Our output is checked
# first, and the input is only consulted when ours fails, to decide whether
# that is a regression we caused or a document the tools reject either way.
run_check() {
	local name="$1" original="$2" output="$3"

	if "${name}_ok" "${output}"; then
		result="ok"
		return
	fi
	if "${name}_ok" "${original}"; then
		result="FAILED"
		regressions+=1
		return
	fi
	result="skip"
	skipped+=1
}

declare -i checked=0 regressions=0 skipped=0
printf '\n%-58s %-8s %s\n' "DOCUMENT" "SCHEMA" "VERIFY"

while read -r original; do
	relative="${original#"${corpus}/"}"
	output="${rendered}/${relative}"
	checked+=1

	if [[ ! -f "${output}" ]]; then
		printf '%-58s %-8s %s\n' "${relative:0:57}" "-" "NOT RENDERED"
		regressions+=1
		continue
	fi

	run_check schema "${original}" "${output}"
	schema_result="${result}"
	run_check verify "${original}" "${output}"
	verify_result="${result}"

	printf '%-58s %-8s %s\n' "${relative:0:57}" "${schema_result}" "${verify_result}"

	if [[ "${schema_result}" == "FAILED" ]]; then
		echo "  the schema rejects our output but accepts the input:"
		report "$(schema_output "${output}")"
	fi
	if [[ "${verify_result}" == "FAILED" ]]; then
		echo "  Verify rejects our output but accepts the input:"
		report "$(verify_output "${output}")"
	fi
done < <(find "${corpus}" -name '*.json' | sort)

# The authored document has no input to compare against: the library built it
# out of Go values, so it simply has to be valid. This is the only check that
# exercises turning Go values into SPDX ones; rendering an existing document
# writes back what it read.
authored="${rendered}/authored.spdx.json"
if [[ -f "${authored}" ]]; then
	echo
	checked+=1
	authored_schema="ok"
	authored_verify="ok"
	if ! schema_ok "${authored}"; then
		authored_schema="FAILED"
		regressions+=1
	fi
	if ! verify_ok "${authored}"; then
		authored_verify="FAILED"
		regressions+=1
	fi
	printf '%-58s %-8s %s\n' "authored.spdx.json (built by the library)" "${authored_schema}" "${authored_verify}"
	if [[ "${authored_schema}" == "FAILED" ]]; then
		echo "  the schema rejects the document the library authored:"
		report "$(schema_output "${authored}")"
	fi
	if [[ "${authored_verify}" == "FAILED" ]]; then
		echo "  Verify rejects the document the library authored:"
		report "$(verify_output "${authored}")"
	fi
else
	echo "the authored document was not written" >&2
	regressions+=1
fi

echo
if [[ ${regressions} -gt 0 ]]; then
	echo "${checked} documents checked, ${regressions} check(s) failed:"
	echo "the SPDX tools reject output we produced from input they accept."
	exit 1
fi
echo "${checked} documents checked. Every document the SPDX tools accept as input, they also accept as our output."
if [[ ${skipped} -gt 0 ]]; then
	echo "Skipped ${skipped} check(s) on input the tools already reject."
fi
