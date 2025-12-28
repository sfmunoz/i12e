#!/usr/bin/env python3
from sys import argv
from base64 import b64decode
from pathlib import Path
from os.path import join
blob_file,ofolder = argv[1],argv[2]
with open(blob_file,"r") as fp:
    buf = fp.read().strip()
for i,line in enumerate(buf.split("\n")):
    ofile,buf = [b64decode(x).decode() for x in line.split("|")]
    p = Path(join(ofolder,ofile))
    p.parent.mkdir(parents=True,exist_ok=True)
    with open(p,"w") as fp:
        fp.write(buf)
