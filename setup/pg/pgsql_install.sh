#!/bin/bash
set -e

sleep 1
psql  -v ON_ERROR_STOP=1 -U postgres  -d remocc -c "CREATE TABLE users(id SERIAL, name VARCHAR(255), email TEXT, orgname VARCHAR(255), groupname VARCHAR(255), password VARCHAR(255), dev TEXT, app TEXT);"
