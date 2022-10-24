#!/usr/bin/env bash
set -euo pipefail
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

xcaddy build --output "$script_dir/caddy" \
  --with github.com/floj/caddy-cors="$script_dir"
