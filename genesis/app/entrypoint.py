import sys
if __name__ == "__main__":
    with open(sys.argv[1],"r") as fp:
        buf = fp.read().strip()
    for i,line in enumerate(buf.split("\n")):
        print(i,line)

