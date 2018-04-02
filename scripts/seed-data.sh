#!/usr/bin/env bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

seedLocal() {
    psql -U alpaca -d alpaca_auth -f ${SCRIPT_DIR}/seed-data-auth.sql
    psql -U alpaca -d alpaca_password_reset -f ${SCRIPT_DIR}/seed-data-password-reset.sql
    psql -U alpaca -d alpaca_role -f ${SCRIPT_DIR}/seed-data-role.sql
    psql -U alpaca -d alpaca_mfa -f ${SCRIPT_DIR}/seed-data-mfa.sql
    psql -U alpaca -d alpaca_email_confirmation -f ${SCRIPT_DIR}/seed-data-email-address-confirmation.sql
}

seedTest() {
    psql -U alpaca -d alpaca_auth_test -f ${SCRIPT_DIR}/seed-data-auth.sql
    psql -U alpaca -d alpaca_password_reset_test -f ${SCRIPT_DIR}/seed-data-password-reset.sql
    psql -U alpaca -d alpaca_role_test -f ${SCRIPT_DIR}/seed-data-role.sql
    psql -U alpaca -d alpaca_mfa_test -f ${SCRIPT_DIR}/seed-data-mfa.sql
    psql -U alpaca -d alpaca_email_confirmation_test -f ${SCRIPT_DIR}/seed-data-email-address-confirmation.sql
}

seedContainer() {
    local path_to_sql_ddl_file=$1
    local container_name=$2
    local db_name=$3
    local db_container_id=$(docker ps -aqf "name=${container_name}")
    echo "Seeding data for container: ${container_name} ${db_container_id} "

    docker cp ${path_to_sql_ddl_file} ${container_name}:/seed-data.sql
    docker exec -i ${db_container_id} /usr/local/bin/psql -d ${db_name} -U alpaca -f /seed-data.sql
}
main() {

    if [[ $1 == "test" ]] ; then
        seedTest

    elif [[ $1 == "local" ]] ; then
        seedLocal

    elif [[ $1 == "docker" ]] ; then
        seedContainer "${SCRIPT_DIR}/seed-data-auth.sql" "alpaca-auth-db" "alpaca_auth"
        seedContainer "${SCRIPT_DIR}/seed-data-password-reset.sql" "alpaca-password-reset-db" "alpaca_password_reset"

    else
        echo "UNHANDLED ARG"
    fi
}

main "$@"