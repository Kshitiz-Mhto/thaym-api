ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'root';
ALTER USER 'root'@'%' REQUIRE NONE;
FLUSH PRIVILEGES;