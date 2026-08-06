#!/usr/bin/env bash
# SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
# SPDX-License-Identifier: Apache-2.0
#
# Syncs the SPDX 3 example documents the round-trip tests check against.
#
# The documents come from the SPDX project's example repository, which
# publishes them under CC0-1.0 (the repository's own code is GPL-3.0-or-later;
# only the documents are vendored here). Run this after upstream publishes new
# examples: it reports what changed and leaves the result staged in the working
# tree for review.
#
#   ./hack/update-spdx-examples.sh
#
# Override the source with SPDX_EXAMPLES_REPO and SPDX_EXAMPLES_REF.

set -euo pipefail

readonly repo_url="${SPDX_EXAMPLES_REPO:-https://github.com/spdx/spdx-examples}"
# Empty means whatever the remote's HEAD points at, so the script keeps
# working if upstream renames its default branch.
readonly repo_ref="${SPDX_EXAMPLES_REF:-}"

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
readonly dest="${root}/testdata/corpus/spdx-examples"
readonly manifest="${root}/testdata/corpus/MANIFEST"

tmp="$(mktemp -d)"
trap 'rm -rf "${tmp}"' EXIT

echo "Cloning ${repo_url} (${repo_ref:-default branch})..."
if [[ -n "${repo_ref}" ]]; then
	git clone --quiet --depth 1 --branch "${repo_ref}" "${repo_url}" "${tmp}/src"
else
	git clone --quiet --depth 1 "${repo_url}" "${tmp}/src"
fi
commit="$(git -C "${tmp}/src" rev-parse HEAD)"
branch="$(git -C "${tmp}/src" rev-parse --abbrev-ref HEAD)"

# An SPDX 3 document is a JSON file bound to a 3.x context and carrying a
# graph. That keeps SPDX 2 examples, JSON schemas and context files out.
find_documents() {
	find . -type f -name '*.json' -printf '%P\n' |
		sort |
		while read -r file; do
			if grep -q 'spdx\.org/rdf/3\.' "${file}" && grep -q '"@graph"' "${file}"; then
				printf '%s\n' "${file}"
			fi
		done
}

mapfile -t documents < <(cd "${tmp}/src" && find_documents)
if [[ ${#documents[@]} -eq 0 ]]; then
	echo "No SPDX 3 documents found in ${repo_url}; refusing to wipe the corpus." >&2
	exit 1
fi

declare -i added=0 updated=0 unchanged=0
for file in "${documents[@]}"; do
	target="${dest}/${file}"
	if [[ ! -f "${target}" ]]; then
		echo "  added:     ${file}"
		added+=1
	elif ! cmp --silent "${tmp}/src/${file}" "${target}"; then
		echo "  updated:   ${file}"
		updated+=1
	else
		unchanged+=1
	fi
	mkdir -p "$(dirname "${target}")"
	cp "${tmp}/src/${file}" "${target}"
done

# Drop anything upstream no longer ships.
declare -i removed=0
if [[ -d "${dest}" ]]; then
	while read -r existing; do
		if ! printf '%s\n' "${documents[@]}" | grep -qxF "${existing}"; then
			echo "  removed:   ${existing}"
			rm "${dest}/${existing}"
			removed+=1
		fi
	done < <(cd "${dest}" && find . -type f -name '*.json' -printf '%P\n' | sort)
	find "${dest}" -type d -empty -delete
fi

cat >"${manifest}" <<EOF
SPDX 3 example documents checked by the round-trip tests.

Source:  ${repo_url}
Ref:     ${repo_ref:-${branch}}
Commit:  ${commit}
License: CC0-1.0, per the source repository's README. Only the SPDX
         documents are vendored, not the repository's build scripts.

Refresh with ./hack/update-spdx-examples.sh. Documents (${#documents[@]}):

$(printf '  %s\n' "${documents[@]}")
EOF

echo
echo "${#documents[@]} documents at ${commit:0:12} (${added} added, ${updated} updated, ${removed} removed, ${unchanged} unchanged)"
echo "Manifest written to ${manifest#"${root}/"}"
echo "Now run: go test ./..."
