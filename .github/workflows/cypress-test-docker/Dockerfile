# Container image that runs your code
FROM cypress/included

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY ../../../frontend /
WORKDIR /frontend

RUN yarn install
RUN yarn cypress run