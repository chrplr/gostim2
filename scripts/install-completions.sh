#!/bin/bash
# install-completions.sh
# Installs bash and zsh completions for gostim2

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BASH_COMPLETION_FILE="$SCRIPT_DIR/gostim2-completion.bash"
ZSH_COMPLETION_FILE="$SCRIPT_DIR/gostim2-completion.zsh"

install_bash() {
    echo "--- Installing Bash Completion ---"
    if [ ! -f "$BASH_COMPLETION_FILE" ]; then
        echo "Error: Bash completion file not found at $BASH_COMPLETION_FILE"
        return 1
    fi

    # Determine destination
    if [ -d /usr/share/bash-completion/completions ]; then
        DEST="/usr/share/bash-completion/completions/gostim2"
    elif [ -d /etc/bash_completion.d ]; then
        DEST="/etc/bash_completion.d/gostim2"
    else
        DEST="$HOME/.local/share/bash-completion/completions/gostim2"
        mkdir -p "$(dirname "$DEST")"
    fi

    echo "Installing to $DEST..."
    if [ -w "$(dirname "$DEST")" ]; then
        cp "$BASH_COMPLETION_FILE" "$DEST"
    else
        sudo cp "$BASH_COMPLETION_FILE" "$DEST"
    fi
    echo "Bash completion installed."
}

install_zsh() {
    echo "--- Installing Zsh Completion ---"
    if [ ! -f "$ZSH_COMPLETION_FILE" ]; then
        echo "Error: Zsh completion file not found at $ZSH_COMPLETION_FILE"
        return 1
    fi

    # Common zsh completion directories
    ZSH_DEST_DIR=""
    if [ -d /usr/local/share/zsh/site-functions ]; then
        ZSH_DEST_DIR="/usr/local/share/zsh/site-functions"
    elif [ -d /usr/share/zsh/site-functions ]; then
        ZSH_DEST_DIR="/usr/share/zsh/site-functions"
    else
        ZSH_DEST_DIR="$HOME/.zsh/completion"
        mkdir -p "$ZSH_DEST_DIR"
        if [[ ":$FPATH:" != *":$ZSH_DEST_DIR:"* ]]; then
            echo "Warning: $ZSH_DEST_DIR is not in your zsh FPATH."
            echo "Add 'fpath=(\$HOME/.zsh/completion \$fpath)' to your .zshrc"
        fi
    fi

    DEST="$ZSH_DEST_DIR/_gostim2"
    echo "Installing to $DEST..."
    if [ -w "$ZSH_DEST_DIR" ]; then
        cp "$ZSH_COMPLETION_FILE" "$DEST"
    else
        sudo cp "$ZSH_COMPLETION_FILE" "$DEST"
    fi
    echo "Zsh completion installed."
}

# Run both
install_bash
echo ""
install_zsh

echo ""
echo "Done. Please restart your shell or source the completion files to enable them."
