#!/bin/sh

# اگر هیچ دستوری داده نشود، help نمایش داده می‌شود
if [ -z "$1" ]; then
    schemagit help
    exit 0
fi

# اجرای دستور SchemaGit
schemagit "$@"
