#!/usr/bin/env bash

recreateUser() {
    local user=$1
    dropuser --if-exists ${user}
    createuser ${user} --createdb
}

dropAlpacaDB() {
    local db_name=$1
    dropdb --if-exists "${db_name}"
    dropdb --if-exists "${db_name}_test"
}

recreateAlpacaDB() {
    local user=$1
    local db_name=$2

    dropAlpacaDB "${db_name}"

    createdb ${db_name} -U ${user}
    createdb ${db_name}_test -U ${user}
}

main() {
    local user="alpaca"

    dropAlpacaDB "alpaca_auth"
    dropAlpacaDB "alpaca_role"
    dropAlpacaDB "alpaca_password_reset"
    dropAlpacaDB "alpaca_mfa"
    dropAlpacaDB "alpaca_email_confirmation"

    recreateUser "$user"

    recreateAlpacaDB "$user" "alpaca_auth"
    recreateAlpacaDB "$user" "alpaca_role"
    recreateAlpacaDB "$user" "alpaca_password_reset"
    recreateAlpacaDB "$user" "alpaca_mfa"
    recreateAlpacaDB "$user" "alpaca_email_confirmation"
}

main "$@"