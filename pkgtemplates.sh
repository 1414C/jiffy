#!/usr/bin/env bash

# pkger is used to bundle the .gotmpl template files into 
# into the binary.  When using parameterized filenames in 
# the codebase, it is neccessary to explicitly tell pkger
# which files to include in pkged.go.

# This script is called from the Makefile to add all of the
# .gotmpl files to the pkged.go source file.  The project 
# will build without the use of the Makefile, but the 
# resulting binary will not run due to the misssing .gotmpl
# files.

# You can always execute the application from the
# source using 'go run main.go ....'.  See the docs at 
# https://1414c.github.io/jiffy for details.

# make sure that pkger is installed
go get -d github.com/markbates/pkger/cmd/pkger

# get a list of all .gotmpl files from the project root
TEMPLATE_LIST=$(find . -name \*.gotmpl)

# convert string to array
TEMPLATE_ARRAY=(${TEMPLATE_LIST})
# echo ${TEMPLATE_ARRAY[@]}

# strip leading '.'' from each array element and insert
# '-include ' for each file to be packaged via pkger.
TEMPLATE_LIST=${TEMPLATE_ARRAY[@]/./'-include '}

# get the unique dirs in TEMPLATE_ARRAY
TEMPLATE_DIRS_ARRAY=''
for i in ${TEMPLATE_ARRAY[@]}; do
    DIR_PATH=$(dirname "$i")
    TEMPLATE_DIRS_ARRAY+=($DIR_PATH)
done
UNIQ_DIRS_ARRAY=($(printf "%s\n" "${TEMPLATE_DIRS_ARRAY[@]}" | sort -u))
#echo "${UNIQ_DIRS_ARRAY[@]}"

# add the UNIQ_DIRS to the TEMPLATE_LIST with the -include
TEMPLATE_LIST+=' '
TEMPLATE_LIST+=${UNIQ_DIRS_ARRAY[@]/./'-include '}

echo '$GOPATH/bin/pkger '$TEMPLATE_LIST
pkger $TEMPLATE_LIST