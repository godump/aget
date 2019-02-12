import subprocess
import sys


def call(command):
    print(command)
    r = subprocess.call(command, shell=True)
    if r != 0:
        sys.exit(r)


def main():
    call('go install github.com/godump/aget')


if __name__ == '__main__':
    main()
