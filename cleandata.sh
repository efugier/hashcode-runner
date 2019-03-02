#!/bin/sh

echo This will erase data/ submissions/ and submissions-tmp/
read -p "Are you sure? [y/n]" -n 1 -r
if [[ $REPLY =~ ^[Yy]$ ]]
then
  rm data/*
  rm submissions/*
  rm submissions-tmp/*
fi
echo ""
