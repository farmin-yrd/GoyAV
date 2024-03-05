#!/bin/bash

dd if=/dev/random of=1M.bin bs=1M count=1
dd if=/dev/random of=2M.bin bs=1M count=2
dd if=/dev/random of=5M.bin bs=1M count=5