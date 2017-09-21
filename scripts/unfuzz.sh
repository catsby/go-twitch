#!/bin/bash

for FILE in `find . -type f` ; do
    # Check if the file contains 'the key'
    if  `grep -q xxxxxxxxxxxxx $FILE`
    then
        if [[ $FILE == *"fuzz.sh"* ]];
        then
          continue
        fi
        echo "Unfuzzing ${FILE}"
        sed -i '' "s/xxxxxxxxxxxxx/${TWITCH_ACCESS_TOKEN}/g" $FILE 
    fi
done
exit 
