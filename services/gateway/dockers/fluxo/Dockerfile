# Use a specific Node.js version
FROM node:20-alpine

# Set working directory
WORKDIR /app

# Install git without version pinning since we're using a specific base image
# Add any other required dependencies in the same layer
RUN apk add --no-cache git

# Clone specific version of repository and install dependencies in one layer
# This reduces the number of layers and image size
RUN git clone --progress https://github.com/PoliGNU-USP/projetoSUS-fluxo.git . && \
  git checkout main && \
  # Install dependencies as a non-root user for better security
  addgroup -S appgroup && \
  adduser -S appuser -G appgroup && \
  chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Install dependencies (including nodemon as a dev dependency if needed)
RUN npm ci  # Using npm ci instead of npm install for more reliable builds

# Expose the port your app runs on (adjust as needed)
EXPOSE 3000

# Health check to verify app is running
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

CMD ["npm", "start"]
