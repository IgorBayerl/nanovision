FROM golang:1.25

ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages for C compilation and Python 3.
RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    xz-utils \
    python3 \
    python3-pip \
    python-is-python3 \
    && rm -rf /var/lib/apt/lists/*

# Install Zig globally
RUN curl -L "https://ziglang.org/download/0.15.2/zig-linux-x86_64-0.15.2.tar.xz" -o zig.tar.xz \
    && tar -xf zig.tar.xz \
    && mv zig-linux-x86_64-0.15.2 /usr/local/zig \
    && rm zig.tar.xz \
    && ln -s /usr/local/zig/zig /usr/local/bin/zig

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