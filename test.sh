TEST_DB_DIALECT=mysql \
TEST_DB_SETUP='mysql -u root -proot -e "create database gotify_test_<time>;"' \
TEST_DB_CONNECTION='root:root@/gotify_test_<time>?charset=utf8&parseTime=True' \
TEST_DB_TEARDOWN='mysql -u root -proot -e "drop database gotify_test_<time>;"' \
go test -v ./...