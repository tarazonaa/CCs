Database
========

The database for the project was created using three distinct technologies:

1. [Postgresql](https://www.postgresql.org/):
 A relational database that helps us keep track of users, images, and our oauth2 implementation
1. [MinIO](https://min.io/)
 The database that helps us store objects for images, to be able to showcase them in history on the frontend.
1. [MongoDB](https://www.mongodb.com/)
 MongoDB is used as a loggin database, this is mostly because logs should be almost untouched.

Env variables
-------------

The env variables for the databases are quite straight forward, such as the root user and passwords, the name of the database, and so on:

```bash
MONGO_INITDB_ROOT_USERNAME=[ROOTUSER]
MONGO_INITDB_ROOT_PASSWORD=[ROOTPWD]
MONGO_INITDB_DATABASE=[ROOTDB]
POSTGRES_DATABASE=[POSTGRESDB]
POSTGRES_PASSWORD=[POSTGRESPWD]
```
