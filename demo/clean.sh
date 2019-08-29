#! /bin/bash

killall -9 basecoin tepleton
TMROOT=./data/chain1/tepleton tepleton unsafe_reset_all
TMROOT=./data/chain2/tepleton tepleton unsafe_reset_all

rm -rf ./data/chain1/basecoin/merkleeyes.db
rm -rf ./data/chain2/basecoin/merkleeyes.db

rm ./*.log

rm ./data/chain1/tepleton/*.bak
rm ./data/chain2/tepleton/*.bak
