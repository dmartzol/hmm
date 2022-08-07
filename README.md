# Hmm! (Hackerspace Membership Management!)

## Description

Fictional project to manage memberships for Makerspaces/Hackerspaces. This is a Work In Progress project combining a backend writen in Golang and a frotend using React. The goal is to include only functionality, contents are not provided in this project.

## Features

* React admin panel to administer resources and sub-resources
* Management of Sessions, cookies, permissions, email verification, accounts, ...
* Management of Hackerspace equipment and equipment authorizations
* PostgreSQL performant database
* Docker Compose for easy bootstrapping

## Requirements

* Docker

## QuickStart

### To run the backend:

```
git clone https://github.com/dmartzol/hmm
cd hmm
make up
make migrate.up
docker compose up --build gateway
```
### Front end is still WIP

## Credits

[![Go](https://www.vectorlogo.zone/logos/golang/golang-ar21.svg)](https://golang.org/ "Golang")
[![PostgreSQL](https://www.vectorlogo.zone/logos/postgresql/postgresql-ar21.svg)](https://www.postgresql.org/ "PostgreSQL")
[![React](https://www.vectorlogo.zone/logos/reactjs/reactjs-ar21.svg)](https://reactjs.org/ "React")
[![Docker](https://www.vectorlogo.zone/logos/docker/docker-ar21.svg)](https://www.docker.com/ "Docker")
[![VectorLogoZone](https://www.vectorlogo.zone/logos/vectorlogozone/vectorlogozone-ar21.svg)](https://www.vectorlogo.zone/ "VectorLogoZone")
