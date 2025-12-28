#!/usr/bin/env python3
from sys import argv
from base64 import b64decode
from pathlib import Path
from os.path import join
blob_file,ofolder = argv[1],argv[2]
with open(blob_file,"r") as fp:
    buf = fp.read().strip()
ofile = ""
for i,line in enumerate(buf.split("\n")):
    entry = b64decode(line).decode()
    if i % 2 == 0:
        ofile = join(ofolder,entry)
        continue
    Path(ofile).parent.mkdir(parents=True,exist_ok=True)
    with open(ofile,"w") as fp:
        fp.write(entry)
