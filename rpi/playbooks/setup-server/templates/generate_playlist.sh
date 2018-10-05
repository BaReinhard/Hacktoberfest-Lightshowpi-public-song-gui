#!/bin/bash

echo "" > /tmp/playlistsongs
cat playlist > /tmp/readinglist
echo "" >> /tmp/readinglist
cat /tmp/readinglist | while IFS=$'\t' read -a line
do
  path_param=${line[1]}
  path=`echo $(eval echo $path_param)`
  Artist=`ffprobe $path 2>&1 | grep -A90 'Metadata:' | grep artist | head -1 | sed 's/.*: //'`
  Title=`ffprobe $path 2>&1 | grep -A90 'Metadata:' | grep title | head -1 | sed 's/.*: //'`
  echo "$Artist ___ $Title" >> /tmp/playlistsongs
done 

cat /tmp/playlistsongs