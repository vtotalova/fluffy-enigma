# Dockerfile for Frontend (React)
# Place this Dockerfile inside the frontend/ directory

# Step 1: Use Node image to build the frontend
FROM node:slim AS build

WORKDIR /app

# Copy package.json and install deps
COPY package.json package-lock.json ./
RUN npm install

# Copy the rest of the frontend source code
COPY . ./

# Build the React app
RUN npm run build

# Step 2: Use nginx to serve the built app
FROM nginx:stable-alpine

# Copy build output to nginx html directory
COPY --from=build /app/build /usr/share/nginx/html

# Expose the port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]