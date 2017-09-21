#!/bin/bash

for FILE in `find . -type f` ; do
    # Check if the file contains 'the key'
    if `grep -q ${TWITCH_ACCESS_TOKEN} $FILE`
    then
        if [[ $FILE == *"fuzz.sh"* ]];
        then
          continue
        fi
        echo "Fuzzing ${FILE}"
        sed -i '' "s/${TWITCH_ACCESS_TOKEN}/xxxxxxxxxxxxx/g" $FILE 
    fi
done
exit 
