#!/usr/bin/env bash

# log=$(git log --oneline | awk '{print $1}')
# for li in $log; do
#     echo "${li},$(git tag -l --points-at $li),$(git log -n 1 --format=%B ${li} --pretty=oneline)"
# done



last_tag=$(git describe --tags `git rev-list --tags --max-count=1`)

commit=`git rev-parse HEAD`
tag_rev=`git rev-parse $last_tag`

if [[ "$commit" == "$tag_rev" ]]
then
    echo "already tagged: $last_tag"
    exit
fi

if [ -z "$last_tag" ]
then
    log=$(git log --pretty=oneline)
    last_tag=0.0.0
else
    log=$(git log $last_tag..HEAD --pretty=oneline)
fi

case "$log" in
    *#major* ) new=$(go run ./cmd/vbum -ver $last_tag -tgt major);;
    *#patch* ) new=$(go run ./cmd/vbum -ver $last_tag -tgt patch);;
    * ) new=$(go run ./cmd/vbum -ver $last_tag -tgt minor);;
esac

echo $new
git tag $new
