import argparse
import osmnx as ox
import os

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-l", "--location", type=str)
    parser.add_argument("-o", "--output_dir", type=str)
    args = parser.parse_args()

    # mkdir
    if not os.path.isdir(args.output_dir):
        os.mkdir(args.output_dir)
    # dl graph
    G = ox.graph_from_place(args.location, network_type="drive", simplify=False)
    # save osmnx stupid repr
    ox.io.save_graphml(G, filepath=os.path.join(args.output_dir, "graph.gml"))

    # save nodes
    idx_to_nodeid = {}
    with open(os.path.join(args.output_dir, "nodes.csv"), "w") as fh:
        for i, node in enumerate(G.nodes(data=True)):
            idx_to_nodeid[i] = node[1]["osmid"]
            # lon lat
            fh.write(f"{node[1]['x']} {node[1]['y']}\n")
    nodeid_to_idx = {v: k for k, v in idx_to_nodeid.items()}

    # save edges
    with open(os.path.join(args.output_dir, "edges.csv"), "w") as fh:
        for edge in G.edges(data=True):
            from_id, to_id, data = edge
            from_idx = nodeid_to_idx[from_id]
            to_idx = nodeid_to_idx[to_id]
            weight = data["length"]

            fh.write(f"{from_idx} {to_idx} {weight}\n")

    # save id mapping
    with open(os.path.join(args.output_dir, "idx_to_nodeid.csv"), "w") as fh:
        for idx, nodeid in idx_to_nodeid.items():
            fh.write(f"{idx} {nodeid}\n")
