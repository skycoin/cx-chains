#!/usr/bin/env bash

# Enable executed commands to be displayed.
set -x

# Ensure GOBIN is set.
if [[ ! -v GOBIN ]]; then
  GOBIN=$(go env GOBIN)
fi

# If CX tracker source dir is defined (in CX_TRACKER_SRC), compile the source.
# The resultant binary should end up in GOBIN.
if [[ -v CX_TRACKER_SRC ]]; then
  ORIG_DIR=${PWD}
  cd "${CX_TRACKER_SRC}" || exit 1
  make install || exit 1
  cd "${ORIG_DIR}" || exit 1
fi

# Compile cxchain.
make install || exit 1

# ENVs.
TRACKER_ADDR=":9091"

# Initiate directory.
function init_temp_dir() {
  TEMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'cx_chains_integration')
  TRACKER_DB="$TEMP_DIR/cx_tracker.db"
}
function clean_temp_dir() {
  rm -rf "$TEMP_DIR"
}

# Start tracker.
function start_tracker() {
  "$GOBIN/cx-tracker" --db="$TRACKER_DB" --addr="$TRACKER_ADDR"
}

init_temp_dir;

start_tracker;

clean_temp_dir;