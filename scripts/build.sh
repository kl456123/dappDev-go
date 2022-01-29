#!/bin/bash

solc --abi contracts/* -o abi/



# weth9
abigen --abi=abi/WETH9.abi --out=generated/weth9/WETH9.go --pkg=weth9
