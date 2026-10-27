#!/usr/bin/env python3
"""Capture what the console actually shows, as a terminal would draw it.

The console uses the alt screen, so a pty capture is a stream of *diffs* against a
buffer that does not exist in the terminal's scrollback. Reading it needs a terminal
emulator — this script uses `pyte` for that.

    pip install pyte
    scripts/console-shot.py --bin ./bin/maitre --cwd /tmp/kb --width 120 --height 40

Two things this gets right that a naive capture does not:

  * the window size is set on the pty *before* the child starts, so the console lays
    out for the size you asked for;
  * the capture is cut at the alt-screen exit, because the exit restores the previous
    screen and replaying past it shows the shell instead of the console.
"""

import argparse
import fcntl
import os
import select
import struct
import sys
import termios
import time

ALT_SCREEN_EXIT = "\x1b[?1049l"


def capture(binary, cwd, keys, cols, rows, settle):
    master, slave = os.openpty()
    fcntl.ioctl(slave, termios.TIOCSWINSZ, struct.pack("HHHH", rows, cols, 0, 0))

    pid = os.fork()
    if pid == 0:
        os.close(master)
        os.setsid()
        fcntl.ioctl(slave, termios.TIOCSCTTY, 0)
        os.dup2(slave, 0)
        os.dup2(slave, 1)
        os.dup2(slave, 2)
        if slave > 2:
            os.close(slave)
        os.chdir(cwd)
        env = dict(os.environ, TERM="xterm-256color")
        os.execve(binary, [binary], env)

    os.close(slave)
    out = bytearray()

    def drain(seconds):
        end = time.time() + seconds
        while time.time() < end:
            ready, _, _ = select.select([master], [], [], 0.1)
            if not ready:
                continue
            try:
                chunk = os.read(master, 65536)
            except OSError:
                return
            if not chunk:
                return
            out.extend(chunk)

    drain(settle)
    for key in keys:
        try:
            os.write(master, key.encode())
        except OSError:
            break
        time.sleep(0.35)
        drain(settle)

    try:
        os.write(master, b"q")
    except OSError:
        pass
    time.sleep(0.3)
    drain(0.5)
    try:
        os.waitpid(pid, 0)
    except ChildProcessError:
        pass
    os.close(master)

    text = out.decode("utf-8", "replace")
    if ALT_SCREEN_EXIT in text:
        text = text.split(ALT_SCREEN_EXIT)[0]
    return text


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--bin", default="./bin/maitre", help="the console binary")
    ap.add_argument("--cwd", default=".", help="where to run it")
    ap.add_argument("--width", type=int, default=120)
    ap.add_argument("--height", type=int, default=40)
    ap.add_argument("--keys", default="", help="comma-separated keys to send first")
    ap.add_argument("--settle", type=float, default=1.0, help="seconds to wait after each key")
    args = ap.parse_args()

    try:
        import pyte
    except ImportError:
        sys.exit("pyte is not installed — `pip install pyte`")

    keys = [k for k in args.keys.split(",") if k]
    raw = capture(args.bin, args.cwd, keys, args.width, args.height, args.settle)

    screen = pyte.Screen(args.width, args.height)
    pyte.Stream(screen).feed(raw)
    for line in screen.display:
        print(line.rstrip())


if __name__ == "__main__":
    main()
