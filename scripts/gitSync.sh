#!/bin/bash

forkDir=~/projects/forks
projectName=$1
user=csun

# functions

function createUpstream() {
    git remote add upstream ${remoteUrl}
}

function prepUpstream() {
    upstreamExists="$(git remote | grep -c upstream)"
    if [ ${upstreamExists} = 0 ];
    then
        echo "upstream does not exist, creating... "
        originUrl=$(git config --get remote.origin.url)
        remoteUrl=$(echo ${originUrl} | sed 's/${user}/${projectName}/g')
        createUpstream
    fi
}

function syncUpstrem() {
    case ${repoDir} in
    *framework|*webapp)
        branchName=ZuoraPanda3
        ;;
    *)
        branchName=master
    esac

    repoName=$(pwd | awk -F "/" '{print $NF}')
    echo ""
    echo "---- syncing up repo ${repoName} branch ${branchName} ----"
    git remote -v
    sleep 2

    echo ""
    echo "---- checkout ----"
    git checkout ${branchName}
    echo ""
    echo "---- fetch upstream ----"
    git fetch upstream
    echo ""
    echo "---- merge upstream ----"
    git merge upstream/${branchName}
    echo ""
    echo "---- push ----"
    git push
}


# main

if [ ! -d ${forkDir}/${projectName} ];
then
    echo "${projectName} does not exist"
    exit 1
fi

for repoDir in ${forkDir}/${projectName}/*;
do
    if [ -d ${repoDir} ];
    then
        cd ${repoDir}
        prepUpstream
    fi
    syncUpstrem
done