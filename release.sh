#!/bin/sh

git tag "v_$1"

git push --tags
git push --all
