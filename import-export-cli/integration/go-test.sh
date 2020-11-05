#!/bin/sh

display_help() {
    echo "Usage: $0 [flags...]" >&2
    echo
    echo "Examples"
    echo "  $0"
    echo "  $0 -run TestVersion"
    echo "  $0 -run .*ApiProduct.*"
    echo "  $0 -v"
    echo "  $0 -logtransport"
    echo
    echo "Flags"
    echo "  -run                Runs only the given test method as a regex"
    echo "  -logtransport       Print http transport request/responses"
    echo "  -v                  Enable verbose logs"
    echo
    exit 1
}

case "$1" in
-h | --help)
  display_help
  exit 0
  ;;
esac

# Checking the OS
case "`uname`" in
Linux*) OS_STR="linux";;
Darwin*) OS_STR="macosx"
esac

# Checking the architecture 32bit/64bit
case "`uname -m`" in
x86_64*) ARCH_STR="x64";;
i386*) ARCH_STR="i586"
esac


ARCHIVE_POSTFIX=$OS_STR-$ARCH_STR.tar.gz
ARCHIVE=$(find ../build/target -name *$ARCHIVE_POSTFIX)
CMD="go test -archive ../$ARCHIVE -p 1 -timeout 0 $*"

echo "Running: $CMD"
eval $CMD
