
Create Database tarea4;
use tarea4;

create table canciones(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` varchar(200),
    album varchar(200),
    `year` varchar(200),
    `rank` varchar(200)
);

select * from canciones;
DELETE from canciones;
