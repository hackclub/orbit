# Orbit

Orbit puts your development environment in the cloud.

- - -

Created at the [LAUNCH Hackathon 2015](http://launch2015.challengepost.com/).
Demo: http://youtu.be/MY01d647S9Y.

## Getting Started

### Prerequisites

##### Client

* Git

##### Server

* Docker (with the client binary)
* Git

### Installation

Use the `orbit` or `orbit-server` command to interact with the app.

    $ go get -u github.com/hackedu/orbit/...
    $ orbit
    $ orbit-server

### Running

##### Server

First, set the `DATABASE_URL` environment variable to a valid Postgres database
(ex: `postgres://orpheus:hunter4@orbit.hackedu.us:5432/orbit`).

Then run these commands to create the DB and serve the backend API.

    # migrate the db
    $ orbit-server createdb
    
    # start the http server
    $ orbit-server -url=http://0.0.0.0:5000 serve

##### Client

Make sure the client's base URL is set to what the server is serving on. Once
the server is up and running, the client should be ready to go.

    $ orbit -h

## License

Copyright (C) 2015 hackEDU

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
