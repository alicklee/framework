#!/bin/bash

#################
################# 如果本仓库 ci 生成 tag 版本失败了，就可以执行这个文件。
#################

tag=$1
if [ "$tag" == "" ]; then
    echo "NEED tag param  =>  v0.x.y "
    exit 1
fi
######################################
# new modify
echo "# " >> build.sh
# new commit
git add .
git commit -m 'Auto tag refresh'

######################################
git tag "${tag}"

######################################
# new modify
echo "# " >> build.sh
# new commit
git add .
git commit -m 'Auto tag refresh'

######################################
# push all
git push && git push --tags
