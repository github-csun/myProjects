#!/bin/bash

branch=$1
release=$2
build=$3

if [ "x${branch}" = "x" ] || [ "x${release}" = "x" ] || [ "x${build}" = "x" ];
  then
    echo "Usage: $0 branch release build"
    echo "branch name (main/release/zdr):"
    read branch
    echo ""
    echo "release name:"
    read release
    echo ""
    echo "build name:"
    read build
    echo ""
fi

if [ ! "${branch}" = "release" ] && [ ! "${branch}" = "main" ] && [ ! "${branch}" = "zdr" ];
  then
    echo "Branch name is not correct, only release or main is accepted"
    echo "Please rerun this script with correct branch name"
    exit 1
fi

resultFolder=testResult
resultFile=${branch}.${release}.${build}.txt

cd ${resultFolder}
if [ -f ${resultFile} ];
  then
    echo "${resultFile} exists, please check and re-try"
    exit 2
fi

vi ${resultFile}

sed -i .bak 's/^ com./com./g' ${resultFile}
mv ${resultFile} ${resultFile}.bak
cat ${resultFile}.bak | awk '{print $1}' > ${resultFile}
rm ${resultFile}.bak

cd ..
ls -l ${resultFolder}/${resultFile}

exit 0
