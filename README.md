# public-song-gui

### Description

This service will be installed on a Raspberry Pi, running as a go server. The server will ping a configured url and post metrics of the state of the lightshow, current song playing, and other songs on the playlist. The server that is hosted on the Public internet will render based on the payload being sent from the RPI. Eventually, functionality will be exposed to have polls to allow others to decide what song will play next, to skip current song, etc.

### Components

- Private Server (Hosted on Raspberry Pi)
- Public Server/Frontend (Hosted on GCP or another service provider, default functionality will rely on GCP)

### Setup

Setup will be available as an Ansible Playbook to configure you're raspberry Pi with the server. There will still require authentication to GCP, but deploy scripts will be made available

### Architecture

**Local Server**

- Runs on Raspberry PI
- Starts with lightshowpi, this will be integrated with the project configs, but scripts will be added to run side by side as well
- Posts Lightshow Pi State to GCP Datastore, this may also be integrated to LightshowPi depending on feasibility

**GCP Server**

- Appengine Standard, keep costs to a minimum (Ideally Free!)
- Serves Static Frontend Files
- Checks for State Changes and Returns them to the Frontend

### Feature Roadmap

Feature Roadmap contains all features that will be targeted, not necessarily in the correct order.

- **_Frontend (COMPLETED)_** - Display Current Song Playing, Songs on Playlist
- **_Frontend_** - Polls for users to cast votes for; next song, skip current song, replay. Possibly more to come
- **_Local Server (COMPLETED)_** - Check Lightshow Pi for current song, next song, playlist songs, if in playlist mode. Send state to Datastore for frontend to render on
- **_Local Server (COMPLETED)_** - Get integrated with the Lightshow Pi Service itself
