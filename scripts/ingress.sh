#!/usr/bin/env bash

main() {
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

    source $SCRIPT_DIR/colors.sh

    local action_verb ingress_url ingress_host
    local services_dir=$SCRIPT_DIR/services
    local action=create

    if [ "$#" -lt "1" ]; then
        echo "Must supply arguments"
        exit 0
    fi

    # idiomatic parameter and option handling in sh
    while test $# -gt 0
    do
        case "$1" in
            --create) action=create; action_verb="creating"
                ;;
            --delete) action=delete; action_verb="deleting"
                ;;
            --*) echo "Unrecognized flag option: $1"
                ;;
            *) echo "Unrecognized argument: $1"
                ;;
        esac
        shift
    done

    print_blue "❤ $action_verb nginx"
    kubectl $action -f $SCRIPT_DIR/ingress/nginx

    if [ "$action" = "create" ] ; then
        print_yellow "❤ waiting for nginx ingress to come alive..."
        ingress_url=$(minikube service nginx-ingress --url)
        print_yellow "❤ nginx ingress is alive at $ingress_url..."
        ingress_host=$(echo $ingress_url | grep -o '[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}')
        print_yellow "❤ adding nginx ingress host to /etc/hosts..."

        echo "$ingress_host    api.alpaca.minikube" | sudo tee -a /etc/hosts
        echo "$ingress_host    alpaca.minikube" | sudo tee -a /etc/hosts

        print_green "/etc/hosts"
        print_blue "$(cat /etc/hosts)"
    elif [ "$action" = "delete" ]
    then
        print_yellow "❤ removing nginx ingress host from /etc/hosts..."
        sudo sed -i '' '/alpaca/d' /etc/hosts
    fi
}

main "$@"
