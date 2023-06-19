#!/bin/sh

$1 -fPIC -shared $2/*.c -I$2 -o $2/built.so 2>&1

