#!/bin/sh
set -eu

REPO="${REPO:-Gladium-AI/x-cli}"
REF="${REF:-main}"
SKILL_NAME="${SKILL_NAME:-x-cli}"
SKILL_PATH="${SKILL_PATH:-skills/${SKILL_NAME}}"
INSTALL_TARGETS="$(printf '%s' "${INSTALL_TARGETS:-claude codex}" | tr ',' ' ')"

download() {
	url="$1"
	out="$2"

	if command -v curl >/dev/null 2>&1; then
		curl -fsSL "$url" -o "$out"
		return
	fi

	if command -v wget >/dev/null 2>&1; then
		wget -q "$url" -O "$out"
		return
	fi

	echo "Error: curl or wget is required" >&2
	exit 1
}

contains_target() {
	target="$1"
	case " $INSTALL_TARGETS " in
		*" $target "*) return 0 ;;
		*) return 1 ;;
	esac
}

resolve_codex_skills_dir() {
	if [ -n "${CODEX_SKILLS_DIR:-}" ]; then
		printf '%s\n' "$CODEX_SKILLS_DIR"
		return
	fi

	if [ -n "${CODEX_HOME:-}" ]; then
		printf '%s/skills\n' "${CODEX_HOME%/}"
		return
	fi

	printf '%s/.codex/skills\n' "$HOME"
}

copy_skill() {
	src="$1"
	skills_dir="$2"
	dst="${skills_dir%/}/${SKILL_NAME}"

	mkdir -p "$skills_dir"
	rm -rf "$dst"
	cp -R "$src" "$dst"
	printf 'Installed skill to %s\n' "$dst"
}

install_skill() {
	src="$1"

	if [ -n "${SKILLS_DIR:-}" ]; then
		copy_skill "$src" "$SKILLS_DIR"
		return
	fi

	if contains_target claude; then
		copy_skill "$src" "${CLAUDE_SKILLS_DIR:-${CLAUDE_CODE_SKILLS_DIR:-$HOME/.claude/skills}}"
	fi

	if contains_target codex; then
		copy_skill "$src" "$(resolve_codex_skills_dir)"
	fi
}

if [ -n "${SOURCE_DIR:-}" ]; then
	src="${SOURCE_DIR%/}/${SKILL_PATH}"
	if [ ! -d "$src" ]; then
		echo "Error: skill directory not found at $src" >&2
		exit 1
	fi
	install_skill "$src"
	exit 0
fi

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT INT TERM

archive="$tmpdir/repo.tar.gz"
url="https://codeload.github.com/${REPO}/tar.gz/refs/heads/${REF}"

download "$url" "$archive"
tar -xzf "$archive" -C "$tmpdir"

root_dir="$(find "$tmpdir" -mindepth 1 -maxdepth 1 -type d | head -n 1)"
if [ -z "$root_dir" ]; then
	echo "Error: could not extract repository archive" >&2
	exit 1
fi

src="${root_dir}/${SKILL_PATH}"
if [ ! -d "$src" ]; then
	echo "Error: skill directory not found in downloaded archive: $src" >&2
	exit 1
fi

install_skill "$src"
