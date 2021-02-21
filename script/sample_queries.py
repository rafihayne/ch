from typing import NamedTuple
import argparse
import os
import math
import random


class Node(NamedTuple):
    idx: int
    long: float
    lat: float


# stupid heuristic. 0.01 seems okay
def euclidean(lhs: Node, rhs: Node) -> float:
    return math.sqrt((lhs.long - rhs.long) ** 2 + (lhs.lat - rhs.lat) ** 2)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-d", "--data_dir", type=str)
    parser.add_argument("-n", "--samples", type=int)
    args = parser.parse_args()

    nodes = []
    with open(os.path.join(args.data_dir, "nodes.csv"), "r") as fh:
        for idx, line in enumerate(fh):
            data = line.strip().split(" ")
            nodes.append(Node(idx, float(data[0]), float(data[1])))

    pairs = []
    # just rejection sample
    while len(pairs) < args.samples:
        start = random.choice(nodes)
        end = start
        while euclidean(start, end) < 0.01:
            end = random.choice(nodes)
        pairs.append((start, end))

    with open(os.path.join(args.data_dir, "samples.csv"), "w") as fh:
        for start, end in pairs:
            fh.write(f"{start.idx} {end.idx}\n")
