# ARG NODE_VERSION=20.9.0

# Use node image for base image for all stages.
# FROM node:${NODE_VERSION}-alpine as base
FROM node:20.9.0 as build


# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the entire application code to the container
COPY . .

# Build the Vue app for production
RUN npm run build

# Use Nginx as the production server
FROM nginx:alpine

# Copy the built Vue app to Nginx's web server directory
COPY --from=build /app/dist/spa /usr/share/nginx/html
COPY --from=build /app/nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80 for the Nginx server
EXPOSE 80

# Start Nginx when the container runs
CMD ["nginx", "-g", "daemon off;"]
