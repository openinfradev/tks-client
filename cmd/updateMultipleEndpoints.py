#!/usr/bin/python3

import argparse
import git
import ruamel.yaml
import os
import sys

print("Entering updateMultipleEndpoint script..")

parser = argparse.ArgumentParser()
parser.add_argument('current_cluster', type=str,
            help="cluster name to which the endpoints are added")
parser.add_argument('endpoint_list', type=str,
            help="endpoint list to add")

args = parser.parse_args()
clusterName = args.current_cluster
endpointListStr = args.endpoint_list

repo = None
config = {}

sitePath = './decapod-site'
siteFileName = "{}/lma/site-values.yaml".format(clusterName)
siteFileNameFull = "{}/{}".format(sitePath, siteFileName)

if not os.path.isdir(sitePath):
    print("Cloning repository...")

    repo = git.Repo.clone_from('https://github.com/robertchoi80/decapod-site', 'decapod-site')
    with repo.config_writer() as git_config:
        git_config.set_value('user', 'email', 'tks-argo@tks.com')
        git_config.set_value('user', 'name', 'TKS Argo')
else:
    repo = git.Repo(sitePath)

with open(siteFileNameFull, 'r') as f:
    config = ruamel.yaml.round_trip_load(f, preserve_quotes=True)

endpointList = endpointListStr.split()
print("endpointList: {}".format(endpointList))

appendCount = 0
charts = config["charts"]
thanosChart = [chart for chart in charts if chart['name'] == "thanos"][0]

for ep in endpointList:
    if (ep in thanosChart['override']['querier.stores']):
        print("The endpoint {} already exists.".format(ep))
    else:
        thanosChart['override']['querier.stores'].append(ep)
        appendCount += 1

if appendCount == 0:
    print("All endpoints already exists. Exiting script..")
    sys.exit(0)

print("After insertion: {}".format(thanosChart))

with open(siteFileNameFull, 'w') as f:
    ruamel.yaml.round_trip_dump(config, f)

diff = repo.git.diff(repo.head.commit.tree)
print(diff)

# Provide a list of the files to stage
repo.index.add([siteFileName])

# Provide a commit message
repo.index.commit("add previous thanos-sidecar endpoints to '{}' cluster".format(clusterName))
repo.remotes.origin.push()

print("Exiting updateMultipleEndpoint script..")
