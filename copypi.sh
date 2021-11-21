#!/bin/sh
if [ $# -lt 1 ]; then
    echo arg error: $*
    exit 1
fi

scp -p co2plotter pi@$1:/home/pi/co2plotter/
