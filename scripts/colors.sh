#!/usr/bin/env bash

red='\e[1;31m'
green='\e[1;32m'
yellow='\e[1;33m'
blue='\e[1;34m'
magenta='\e[1;35m'

close='\e[0m'

function print_color {
    printf "$1%s$close\n" "$2"
}

function print_red {
    print_color "$red" "$1"
}

function print_green {
    print_color "$green" "$1"
}

function print_yellow {
    print_color "$yellow" "$1"
}

function print_blue {
    print_color "$blue" "$1"
}

function print_magenta {
    print_color "$magenta" "$1"
}