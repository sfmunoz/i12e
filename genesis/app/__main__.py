#!/usr/bin/env python3
import sys
from os import getenv
from logging import getLogger, basicConfig, INFO
from .butane import Butane
from .artifact import Artifact

basicConfig(format='%(asctime)s [%(relativeCreated)7.0f] [%(levelname).1s] %(message)s (%(module)s:%(lineno)d)',level=INFO,stream=sys.stderr)
log = getLogger(__name__)

def main():
    if getenv("GENESIS_ARTIFACT") == "1":
        Artifact().run()
        sys.exit(0)
    genesis_target = getenv("GENESIS_TARGET")
    if genesis_target is not None and genesis_target != "":
        Butane(genesis_target).run()
        sys.exit(0)

if __name__ == "__main__":
    main()
