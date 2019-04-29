#!/bin/bash

# BASE_DIR is the project root absolute path
export BASE_DIR=$PWD/
# DB Test Exports
export POSTGRES_PORT=5432
export POSTGRES_HOST=localhost
export POSTGRES_USER=postgres
export POSTGRES_DB=remocc
export POSTGRES_PASSWORD=remocc
# SSH Test Exports
export SSHD_CMD=/usr/sbin/sshd
export AUTHORIZED_KEYS=$BASE_DIR/test/sshtests/authorized_keys
