#!/bin/bash
set -e

function undpack_and_install_tool {
	TOOL_NAME=$1
	TOOL_URL=$2
	if [ ! -f "$HOME/bin/$TOOL_NAME" ]; then
		echo "installing [$TOOL_NAME] ..."
		wget "$TOOL_URL" -O "$HOME/bin/$TOOL_NAME.tgz"
		tar -xvzf "$HOME/bin/$TOOL_NAME.tgz" -C "$HOME/bin/"
		rm -f "$HOME/bin/$TOOL_NAME.tgz"
		chmod +x "$HOME/bin/$TOOL_NAME"
	fi
}

# $HOME/bin
if [ ! -d "$HOME/bin" ]; then mkdir "$HOME/bin"; fi
export PATH="$HOME/bin:$PATH"

# install tools
undpack_and_install_tool "pack" "https://github.com/buildpacks/pack/releases/download/v0.11.2/pack-v0.11.2-linux.tgz"
