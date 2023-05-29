#!/bin/sh

gcc -fPIC -shared $1/*.c -I$1 -o $1/built.so 2>&1

