#!/bin/bash

fyne package -os wasm
rm -r serve/wasm
mv wasm serve/
