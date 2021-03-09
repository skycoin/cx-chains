#!/usr/bin/env bash

# Enable executed commands to be displayed.
#set -x

# Initiate environment.
function init_env() {
  echo "<< Initiating environment. >>"

  # Ensure GOBIN is set.
  : "${GOBIN:="$HOME/go/bin"}"
  echo "GOBIN=$GOBIN"

  # Ensure TEMP_DIR is set.
  : "${TEMP_DIR:=$(mktemp -d 2>/dev/null || mktemp -d -t 'cx_chains_integration')}"
  echo "TEMP_DIR=$TEMP_DIR"

  # Ensure TRACKER_SRC is set.
  : "${TRACKER_SRC=$(dirname "$(pwd)")/cx-tracker}"
  echo "TRACKER_SRC=$TRACKER_SRC"

  # Ensure TRACKER_ADDR is set.
  : "${TRACKER_ADDR:=":9091"}"
  echo "TRACKER_ADDR=$TRACKER_ADDR"

  # Ensure TRACKER_DB is set.
  : "${TRACKER_DB:="$TEMP_DIR/cx_tracker.db"}"
  echo "TRACKER_DB=$TRACKER_DB"

  # Ensure TRACKER_LOG is set.
  : "${TRACKER_LOG:="$TEMP_DIR/cx_tracker.log"}"
  echo "TRACKER_LOG=$TRACKER_LOG"

  # Ensure TRACKER_PID is set.
  : "${TRACKER_PID:="$TEMP_DIR/cx_tracker.pid"}"
  echo "TRACKER_PID=$TRACKER_PID"

  # Ensure CXCHAIN_DIR is set.
  : "${CXCHAIN_DIR:="$TEMP_DIR/cxchain"}"

}

# Initiate binaries.
function init_bin() {
  echo "<< Initiating binaries. >>"

  # Compile cx-tracker.
  _d=$(pwd)
  cd "$TRACKER_SRC" || exit 1
  echo ">> Installing 'cx-tracker'."
  make install || exit 1
  cd "$_d" || exit 1

  # Compile cxchain.
  echo ">> Installing 'cxchain'."
  make install || exit 1
}

# Clean temp dir.
function clean_temp_dir() {
  echo "<< Cleaning temporary directory. >>"
  rm -rf "$TEMP_DIR"
}

# Start tracker.
function start_tracker() {
  echo "<< Starting 'cx-tracker'. >>"
  "$GOBIN/cx-tracker" --db="$TRACKER_DB" --addr="$TRACKER_ADDR" >> "$TRACKER_LOG" 2>&1 &
  echo $$ > "$TRACKER_PID"
}

# Stop tracker.
function stop_tracker() {
  echo "<< Stopping 'cx-tracker'. >>"
  cat "$TRACKER_PID" || xargs kill
}

# Start cxchain client.
function start_cxchain_client() {
  echo "<< Starting 'cxchain'. >>"
  if [[ "$?" -ne 3 ]]; then echo "Needs 3 arguments." && exit 1; fi

  local _index="$1"
  local _port="$2"
  local _web_port="$3"

  "$GOBIN/cxchain" --enable-all-api-sets --client \
    --data-dir="$CXCHAIN_DIR/$_index" \
    --port="$_port" \
    --web-interface-port="$_web_port" \
    >> "$TEMP_DIR/cxchain_$_index.log" 2>&1 &

  echo $$ > "$TEMP_DIR/cxchain$_index.pid"
}

# Start cxchain master.
function start_cxchain_master() {
  echo "<< Starting 'cxchain'. >>"
  if [[ "$?" -ne 3 ]]; then echo "Needs 3 arguments." && exit 1; fi

  local _index="$1"
  local _port="$2"
  local _web_port="$3"

  "$GOBIN/cxchain" --enable-all-api-sets \
    --data-dir="$CXCHAIN_DIR/$_index" \
    --port="$_port" \
    --web-interface-port="$_web_port" \
    >> "$TEMP_DIR/cxchain_$_index.log" 2>&1 &

  echo $$ > "$TEMP_DIR/cxchain$_index.pid"
}

init_env;
init_bin;

start_tracker;

sleep 5

stop_tracker;

clean_temp_dir;