### First stage: build backend
FROM golang:1.14-alpine AS go
# Set the working directory
WORKDIR /src
# Download dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download
# Import code from the context
COPY . .
# Build the go binary
RUN CGO_ENABLED=0 go build -o ./ft ./cmd/feature-toggle/.

### Second stage: build frontend
FROM node:14-buster AS node
WORKDIR /src
# Import the frontend code from the context
COPY ./static/package.json ./static/yarn.lock ./
# Install all dependencies
RUN yarn install
# Copy source files from the context
COPY ./static ./
# Run a script target to build a project in production mode
RUN npm run lint && npm run prod

### Last stage: build the app
FROM scratch
# Set the working directory
WORKDIR /opt/ft
# Further actions will be performed as a non-privileged user
USER nobody:nogroup
# Copy everything required from the go stage into the final stage
COPY --from=go /etc/passwd /etc/passwd
COPY --from=go /etc/group /etc/group
# Copy everything required from the go stage into the final stage
COPY --from=go --chown=nobody:nogroup /src/ft ./ft
# Copy the dist directory from the node stage
COPY --from=node --chown=nobody:nogroup /src/dist/ft ./static/dist/ft
# Install required packages
EXPOSE 8080
# Further actions will be performed as a non-privileged user
USER nobody:nogroup
ENTRYPOINT ["./ft"]
