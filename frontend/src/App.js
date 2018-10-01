import React, { Component } from "react";
import logo from "./logo.svg";
import axios from "axios";
import "./App.css";

class App extends Component {
  constructor() {
    this.state = {
      psgState: {}
    };
  }
  async componentDidMount() {
    let newState = {};
    if (process.env.PSG_ENV === "PROD") {
      newState = await axios.get("/api/getState").then(res => {
        res;
        if (res.data) {
          return res.data;
        }
      });
    } else {
      newState = await new Promise((resolve, reject) => {
        return {
          songs: [
            { name: "First Song", artist: "First Artist" },
            { name: "Second Song", artist: "Second Artist" },
            { name: "Third Song", artist: "Third Artist" }
          ],
          status: "playing",
          currentSongIndex: 0,
          currentSong: { name: "First Song", artist: "First Artist" }
        };
      });
    }
    this.setState({ psgState: newState });
  }
  render() {
    if (this.state.psgState) {
      return (
        <div className="App">
          <header className="App-header">
            <img src={logo} className="App-logo" alt="logo" />
            <h1 className="App-title">Welcome to React</h1>
          </header>
          <p className="App-intro">
            To get started, edit <code>src/App.js</code> and save to reload.
          </p>
        </div>
      );
    } else {
      return (
        <div className="App">
          <header className="App-header">
            <img src={logo} className="App-logo" alt="logo" />
            <h1 className="App-title">Welcome to React</h1>
          </header>
          <div>
            <h1>Current Song</h1>
            <div>
              <div>
                <h3>Artist</h3>
                {this.state.psgState.currentSong.artist}
              </div>
              <div>
                <h3>Song</h3>
                {this.state.psgState.currentSong.name}
              </div>
            </div>
          </div>
          <div>
            <h1>Songs in Playlist</h1>
            {this.state.psgState.songs.map(song => (
              <div>
                <div>
                  <h3>Artist</h3>
                  {song.artist}
                </div>
                <div>
                  <h3>Song</h3>
                  {song.name}
                </div>
              </div>
            ))}
          </div>
        </div>
      );
    }
  }
}

export default App;
