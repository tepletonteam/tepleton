#! /bin/bash

killall -9 basecoin tepleton
TMHOME=./data/chain1 tepleton unsafe_reset_all
TMHOME=./data/chain2 tepleton unsafe_reset_all

rm ./*.log

rm ./data/chain1/*.bak
rm ./data/chain2/*.bak
