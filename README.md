# Web Radio Share (WRS)

Little web server to share the song your are listening

## Why this project ?

Often at school we (students) are listening music. And friens ask what are we
currently listening and what type of music we listen.

This project is to share what we currently listening and show friends what
songs and music we listen.

## What is this project ?

This is a basic web music player.

Other users can come on the web page and listen what the hoster listen or they
can listen others songs.

When the hoster start the server and go the web page at the address
[localhost:5596/hoster](http://localhost:5596/hoster "The web page of the running project").  
In the web page he can see all the songs he have in a choosen directory
(default is `~/Music`).
He must be logged in to change the sound (so no other users can change it).

Other users have to go at the address of the host.  
`<host_address>:5596/`  
They see what is in the music library of the hoster, what song he currently
listening and lister that one at the same time as the hoster.

## How can I install this project and share my sounds

1. Clone or download this repository.
2. Be sure you have docker installed otherwise install it ([Docker install](https://docs.docker.com/get-docker/ "Install docker"))
3. Complete the `env.json.tpl` file and rename it to `env.json`
3. Start the container by doing `docker-compose up`
4. The project is started.

## How can I configure the project ?

The project can be configured throught the `env.json` file.  
You have the `env.json.tpl` file to help make your own configuration.
