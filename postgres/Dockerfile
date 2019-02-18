FROM postgres:10.3

COPY 1_initial.up.sql /docker-entrypoint-initdb.d
COPY 2_test_users.sql /docker-entrypoint-initdb.d
COPY 3_test_contacts.sql /docker-entrypoint-initdb.d
COPY 4_test_dms.sql /docker-entrypoint-initdb.d
COPY 5_test_group.sql /docker-entrypoint-initdb.d
