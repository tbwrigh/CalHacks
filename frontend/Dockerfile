# Use a base image with Node.js
FROM node:18

# Set the working directory
WORKDIR /app

# Copy the package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application files
COPY . .

# Build the SvelteKit app
# RUN npm run build

EXPOSE 5173

# Start the SvelteKit server
CMD ["npm", "run", "dev", "--", "--host"]
