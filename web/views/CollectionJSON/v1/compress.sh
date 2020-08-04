#!/bin/bash

rm -f *.min.json

for file_json in *.json
do
	file_min_json=${file_json%.*}.min.${file_json#*.}
	cat $file_json | tr -d '\n\t' | tr -s ' ' | sed 's:/\*.*\*/::g' | sed 's/ \?\([{}();,:]\) \?/\1/g' > $file_min_json
done
