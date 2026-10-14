#!/bin/bash

TMP_FILE="/tmp/ibgateway"
DEPS_PATH="$HOME/deps"
IBGATEWAY_PATH="$DEPS_PATH/ibgateway"
# From https://www.interactivebrokers.com/campus/ibkr-quant-news/interactive-brokers-gateway-install-setup/
IB_DOWNLOAD_URL="https://download.interactivebrokers.com/installers/ibgateway/ibgateway-latest-standalone-linux-x64.sh"

usage() {
  echo "usage: $0 [--no-clean]"
}

CLEAN=1
if [ "$#" -gt 0 ]; then
  if [ "$1" = "--no-clean" ]; then
    CLEAN=0 
  else
    usage
    exit 1
  fi
fi

# TODO: Check the installed vs latest version?

if [ "$CLEAN" -eq 1 ]; then
  echo "Cleaning up"
  rm -rf "$IBGATEWAY_PATH"
fi

echo "Reminder: install IBGateway in $DEPS_PATH.

It will automatically create an $IBGATEWAY_PATH/{version} directory"

exit 0
curl -o "$TMP_FILE" "$IB_DOWNLOAD_URL"
chmod u+x "$TMP_FILE"
"$TMP_FILE"

printf 'Clean install file (Y/n)? '
read answer

if [ "$answer" = "" -o "$answer" != "${answer#[Yy]}" ]; then 
  rm "$TMP_FILE"
else
  echo "Left $TMP_FILE..."
fi
