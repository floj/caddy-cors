#!/usr/bin/env bash
set -euo pipefail
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"


caddy_version=2.4.6

xcaddy build "v$caddy_version" --output "$script_dir/caddy" \
  --with github.com/floj/caddy-cors="$script_dir"
