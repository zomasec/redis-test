# Use Node.js base image
FROM node:14

# Install necessary dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    telnet \
    bash \
    && rm -rf /var/lib/apt/lists/*
# Set the working directory
WORKDIR /app

# Copy package.json and install dependencies
COPY package.json /app/
RUN npm install

# Copy the app code
COPY . /app

# Expose port 3000 (optional, if you serve any API in future)
EXPOSE 3000

# Run the application
CMD ["npm", "start"]
