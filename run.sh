#!/bin/bash

workdir=$(cd $(dirname $0); pwd)
"${workdir}"/bin/myProject -f src/conf/app.ini # -z ${PWD}/log