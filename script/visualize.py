import matplotlib

matplotlib.use("Agg")

import osmnx as ox
import argparse
import os

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-d", "--data_dir", type=str)
    args = parser.parse_args()

    G = ox.load_graphml(os.path.join(args.data_dir, "graph.gml"))

    idx_to_nodeid = {}
    with open(os.path.join(args.data_dir, "idx_to_nodeid.csv"), "r") as fh:
        for line in fh:
            idx, id_ = line.strip().split()
            idx_to_nodeid[int(idx)] = int(id_)
    nodeid_to_idx = {v: k for k, v in idx_to_nodeid.items()}

    paths = []
    with open(os.path.join(args.data_dir, "paths.csv"), "r") as fh:
        for i, line in enumerate(fh):
            path = []
            for idx in line.strip().split(","):
                if idx:
                    path.append(idx_to_nodeid[int(idx)])
            if path:
                paths.append(path)

    fig, ax = ox.plot_graph_routes(
        G, paths, figsize=(35, 35), node_size=0.0, edge_color="#01B9FF", bgcolor="white"
    )
    fig.savefig(os.path.join(args.data_dir, "fig.png"))
