#!/bin/bash

set -x
set -e

rm -rf etcd-configs && mkdir etcd-configs
./declcfg blob olm.package,olm.bundle quay.io/openshift-community-operators/etcd:v0.9.4 >> etcd-configs/out
./declcfg blob olm.bundle quay.io/openshift-community-operators/etcd:v0.9.2 >> etcd-configs/out
./declcfg blob olm.bundle quay.io/openshift-community-operators/etcd:v0.9.0 >> etcd-configs/out
./declcfg-inline-bundles etcd-configs  --delete-non-head-objects
mv etcd-configs/etcd/etcd.json etcd-configs/out
mv etcd-configs/etcd/objects/etcdoperator.v0.9.0/etcdbackups.etcd.database.coreos.com.crd.yaml etcd-configs/etcdbackups.etcd.database.coreos.com.crd.0.9.0.yaml
mv etcd-configs/etcd/objects/etcdoperator.v0.9.0/etcdclusters.etcd.database.coreos.com.crd.yaml etcd-configs/etcdclusters.etcd.database.coreos.com.crd.0.9.0.yaml
mv etcd-configs/etcd/objects/etcdoperator.v0.9.0/etcdrestores.etcd.database.coreos.com.crd.yaml etcd-configs/etcdrestores.etcd.database.coreos.com.crd.0.9.0.yaml
mv etcd-configs/etcd/objects/etcdoperator.v0.9.0/etcdoperator.v0.9.0.clusterserviceversion.yaml etcd-configs/etcdoperator.v0.9.0.clusterserviceversion.yaml
mv etcd-configs/etcd/objects/etcdoperator.v0.9.4/etcdbackups.etcd.database.coreos.com.crd.yaml etcd-configs/etcdbackups.etcd.database.coreos.com.crd.0.9.4.yaml
mv etcd-configs/etcd/objects/etcdoperator.v0.9.4/etcdclusters.etcd.database.coreos.com.crd.yaml etcd-configs/etcdclusters.etcd.database.coreos.com.crd.0.9.4.yaml
mv etcd-configs/etcd/objects/etcdoperator.v0.9.4/etcdrestores.etcd.database.coreos.com.crd.yaml etcd-configs/etcdrestores.etcd.database.coreos.com.crd.0.9.4.yaml
mv etcd-configs/etcd/objects/etcdoperator.v0.9.4/etcdoperator.v0.9.4.clusterserviceversion.yaml etcd-configs/etcdoperator.v0.9.4.clusterserviceversion.yaml
sed -i 's|objects/etcdoperator.v0.9.0/etcdbackups.etcd.database.coreos.com.crd.yaml|etcdbackups.etcd.database.coreos.com.crd.0.9.0.yaml|' etcd-configs/out
sed -i 's|objects/etcdoperator.v0.9.0/etcdclusters.etcd.database.coreos.com.crd.yaml|etcdclusters.etcd.database.coreos.com.crd.0.9.0.yaml|' etcd-configs/out
sed -i 's|objects/etcdoperator.v0.9.0/etcdrestores.etcd.database.coreos.com.crd.yaml|etcdrestores.etcd.database.coreos.com.crd.0.9.0.yaml|' etcd-configs/out
sed -i 's|objects/etcdoperator.v0.9.0/etcdoperator.v0.9.0.clusterserviceversion.yaml|etcdoperator.v0.9.0.clusterserviceversion.yaml|' etcd-configs/out
sed -i 's|objects/etcdoperator.v0.9.4/etcdbackups.etcd.database.coreos.com.crd.yaml|etcdbackups.etcd.database.coreos.com.crd.0.9.4.yaml|' etcd-configs/out
sed -i 's|objects/etcdoperator.v0.9.4/etcdclusters.etcd.database.coreos.com.crd.yaml|etcdclusters.etcd.database.coreos.com.crd.0.9.4.yaml|' etcd-configs/out
sed -i 's|objects/etcdoperator.v0.9.4/etcdrestores.etcd.database.coreos.com.crd.yaml|etcdrestores.etcd.database.coreos.com.crd.0.9.4.yaml|' etcd-configs/out
sed -i 's|objects/etcdoperator.v0.9.4/etcdoperator.v0.9.4.clusterserviceversion.yaml|etcdoperator.v0.9.4.clusterserviceversion.yaml|' etcd-configs/out
(mkdir -p etcd-configs/something/random && cd etcd-configs/something/random && wget https://gist.githubusercontent.com/joelanford/33cc49a44ff214987ec1aa17c440e6f8/raw/084464d793d7aa0c98aa462077b066b274627c3a/gistfile1.txt)
./declcfg validate etcd-configs
(cd etcd-configs && git init && git add . && git commit -m "initial commit")
