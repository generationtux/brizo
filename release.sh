#!/usr/bin/env sh
# release.sh is a very simple release script for Github repos

if [ "$1" == "--help" ] || [ "$1" == "-h" ] || [ "$1" == "" ]; then
    echo "Usage: `basename $0` <owner> <repository> <version>"
    exit 0
fi

if [ ! -r .github_token ]; then
    echo 'Missing .github_token file to use for authentication.\nPlease add "ACCESS_TOKEN=:yourtoken" to a .github_token file.'
    exit 0
fi

owner=$1
repository=$2
version=$3

github_token=$(cat .github_token | sed -e 's/ACCESS_TOKEN=//')
payload=$(printf '{"tag_name": "v%s","target_commitish": "master","name": "v%s","body": "Release of version %s","draft": false,"prerelease": false}' $version $version $version)
url=$(printf 'https://api.github.com/repos/%s/%s/releases?access_token=%s' $owner $repository $github_token)

curl --data "$payload" $url
