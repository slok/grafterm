#!/bin/bash

current_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [[ -n ${CIRCLE_TAG} ]]; then
    echo "Tag ${CIRCLE_TAG}. building releases..."

    archs=( Linux Darwin Windows ARM64 ARM )
    for arch in "${archs[@]}"
    do
        VERSION=${CIRCLE_TAG} ostype=${arch} ${current_dir}/build.sh
    done
else
    echo "no tag, skipping release..."
fi