# Use a small, production-ready web server image (Nginx Alpine)
FROM nginx:alpine

# Remove default nginx static files
RUN rm -rf /usr/share/nginx/html/*

# Copy frontend files to nginx html directory
COPY index.html /usr/share/nginx/html/
COPY app.js /usr/share/nginx/html/

# Expose port 80 (default for nginx)
EXPOSE 80

# Nginx runs as non-root by default in this image
CMD ["nginx", "-g", "daemon off;"]