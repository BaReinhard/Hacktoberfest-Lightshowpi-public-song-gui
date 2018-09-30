# public-song-gui

***Description***
This service will be installed on a Raspberry Pi, running as a go server. The server will ping a configured url and post metrics of the state of the lightshow, current song playing, and other songs on the playlist. The server that is hosted on the Public internet will render based on the payload being sent from the RPI. Eventually, functionality will be exposed to have polls to allow others to decide what song will play next, to skip current song, etc.

***Components***
- Private Server (Hosted on Raspberry Pi)
- Public Server/Frontend (Hosted on GCP or another service provider, default functionality will rely on GCP)

***Setup***
Setup will be available as an Ansible Playbook to configure you're raspberry Pi with the server. There will still require authentication to GCP, but deploy scripts will be made available
