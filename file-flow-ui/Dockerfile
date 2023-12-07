# Stage 1: Compile and Build angular codebase

#  Use office node image as the base image
FROM node:latest as build

# Set the working directory
WORKDIR /app

# Add the source code to app
COPY . .

# Install dependencies
RUN npm install

# Generate the build of the applicaion
RUN npm run build-dev

# Stage 2: Serve app with nginx server

# Use official nginx image as the base image
FROM nginx:latest

# COPY the build output to replace the deault nginx contents.
COPY --from=build /app/dist/file-flow-ui /usr/share/nginx/html

# Expose port 80
EXPOSE 80
