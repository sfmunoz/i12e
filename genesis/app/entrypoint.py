#!/usr/bin/env python3
from base64 import b64decode
from pathlib import Path
blob_file = Path(__file__).resolve().parent / "code.blob"
ofolder = Path("/app")
for line in blob_file.read_text().strip().split("\n"):
    ofile,buf = [b64decode(x).decode() for x in line.split("|")]
    p = ofolder / ofile
    p.parent.mkdir(parents=True,exist_ok=True)
    p.write_text(buf)
