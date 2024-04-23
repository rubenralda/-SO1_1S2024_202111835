# Dockerfile
Iniciar sesión antes de docker hub con:
~~~
docker login -u nombreusuario
~~~
y pedira la contraseña la cual se puede obtener un token de sesion en la pagina y nuestro perfil.

Para hacer build de nuestro archivo dockerfile del client es:
~~~
docker build --tag rubenralda/consumer .
~~~
y subir:
~~~
docker push rubenralda/consumer:latest
~~~