#!/bin/sh

echo "[]" > ./imagestore_info.json

cd registry
rm -rf blobs/*
rm -rf blob-manifests/*
rm -rf repos/*
rm -rf temp/*

