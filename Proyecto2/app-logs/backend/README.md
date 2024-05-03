# Dependencias
~~~sh
//comando utilizados
npm install express
npm install morgan
npm install nodemon -D
npm install cors
npm install mongodb@6.5
~~~
# Ejecutar proyecto
Puerto configurado 3000
~~~
npm run dev
~~~

# Endpoint
~~~
http://{{host}}/api/logs
~~~

# Docker
sudo docker build -t rubenralda/backend-sopes-pj2:1.0 .
sudo docker push rubenralda/backend-sopes-pj2:1.0