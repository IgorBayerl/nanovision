FROM golang:1.25

ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages for C compilation, and Windows cross-compilation
# Also install Python 3. These commands are run as root.
RUN apt-get update && apt-get install -y \
    build-essential \
    mingw-w64 \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user to avoid running as root
RUN useradd -ms /bin/bash vscode

# Create the GOPATH directory structure and change its ownership to the new user.
# This MUST be done as root BEFORE switching the user.
RUN mkdir -p /go/src /go/bin && chown -R vscode:vscode /go

# Now, switch to the non-root user for all subsequent commands
USER vscode

# The base golang image already sets GOPATH=/go and adds it to the PATH,
# so we don't need to explicitly set those ENV variables again.

# Set the working directory for the user's shell
WORKDIR /workspaces