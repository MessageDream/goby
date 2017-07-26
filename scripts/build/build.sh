#!/bin/bash

usage="Usage: $0 -[v] Version"

version=$1

workingdir=$(pwd) 

src_root_dir=${workingdir}"/../../"
src_main_file=${src_root_dir}"/goby.go"
exc_main_file=${workingdir}"/goby"
exc_windows_main_file=${exc_main_file}".exe"

src_conf_dir=${src_root_dir}"/conf"
src_public_dir=${src_root_dir}"/public"
src_template_dir=${src_root_dir}"/template"

src_scripts_dir=${src_root_dir}"/scripts"
src_scripts_sys_dir=${src_scripts_dir}"/systemd"
src_scripts_sql_dir=${src_scripts_dir}"/sql"


function copy_exc(){
    if [ "$1"x = "windows"x ];then
        $(cp ${exc_windows_main_file} $2) 
        $(rm ${exc_windows_main_file})
    else
        $(cp ${exc_main_file} $2) 
        $(rm ${exc_main_file})
    fi
}

function copy_res(){
    $(mkdir $1)
    os_build_dir=$1"/" 
    $(cp -rf ${src_conf_dir} ${os_build_dir})
    $(cp -rf ${src_public_dir} ${os_build_dir})
    $(cp -rf ${src_template_dir} ${os_build_dir})
    $(cp -rf ${src_scripts_sys_dir} ${os_build_dir}scripts/)
    $(cp -rf ${src_scripts_sql_dir} ${os_build_dir}scripts/)
    echo ${os_build_dir}
}

function build_zip(){
   
    $(CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build ${src_main_file})
    
    copy_exc $1 $3

    name="goby-v${version}-$1-$2"
    zip_name=${name}".zip"

    $(mv $3 ${name})

    $(zip -rq ${zip_name} ${name})

    $(mv ${name} $3)

}

function start() {
    echo $1
    os_dir=$(copy_res $1)
    build_zip $1 "amd64" ${os_dir}
    build_zip $1 "386" ${os_dir}

    $(rm -rf ${os_dir})
}

start "linux" 
start "darwin"
start "windows"


