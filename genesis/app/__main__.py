#!/usr/bin/env python3
import sys
from os import execlp
from logging import getLogger, basicConfig, INFO
from argparse import ArgumentParser, RawTextHelpFormatter
from .butane import Butane
from .artifact import Artifact

basicConfig(format='%(asctime)s [%(relativeCreated)7.0f] [%(levelname).1s] %(message)s (%(module)s:%(lineno)d)',level=INFO,stream=sys.stderr)
log = getLogger(__name__)

BUTANE_K3S_VERSION_DEFAULT = "v1.34.3+k3s1"
BUTANE_VALID_OUTPUTS = ["bash_b64","bash_raw","ignition","debug"]
BUTANE_DEFAULT_OUTPUT = "bash_b64"

def genesis_run(args):
    if args.command == 'artifact':
        Artifact().run()
        return
    elif args.command == 'butane':
        Butane(args).run()
        return
    elif args.command in ['python3','sh']:
        execlp(args.command,args.command)
    raise Exception("unknown command '{0}'".format(args.command))

def main():
    if len(sys.argv) < 2:
        sys.argv.append("-h")  # enforce help display with no arguments
    epilog = "46285520+sfmunoz@users.noreply.github.com (C) 2026"
    parser = ArgumentParser(
        description = 'genesis',
        epilog = epilog,
        formatter_class = RawTextHelpFormatter,
    )
    parser.add_argument('-d', '--debug', action='store_true',
                        help='enable debug mode')
    subparsers = parser.add_subparsers(
        title = 'genesis command',
        description = 'choose one genesis command',
        help = 'genesis command to be run',
        dest = 'command',
    )
    parser_artifact = subparsers.add_parser('artifact', help='generate artifact and push it using rclone')
    parser_artifact.set_defaults(func=genesis_run)
    parser_butane = subparsers.add_parser('butane', help='run butane to generate ignition code')
    parser_butane.add_argument('-o', '--output', metavar='output', action='store',
                        dest='output', type=str, choices=BUTANE_VALID_OUTPUTS, default=BUTANE_DEFAULT_OUTPUT,
                        help='output (default: {0}; valid: {1})'.format(BUTANE_DEFAULT_OUTPUT,", ".join(BUTANE_VALID_OUTPUTS)))
    parser_butane.set_defaults(func=genesis_run)
    parser_python3 = subparsers.add_parser('python3', help='run python3 within the container')
    parser_python3.set_defaults(func=genesis_run)
    parser_sh = subparsers.add_parser('sh', help='run sh within the container')
    parser_sh.set_defaults(func=genesis_run)
    args = parser.parse_args()
    if args.debug:
        from logging import DEBUG
        log.setLevel(DEBUG)
    args.func(args)

if __name__ == "__main__":
    main()
