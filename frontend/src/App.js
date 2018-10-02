import React, { Component } from "react";
import logo from "./logo.svg";
import axios from "axios";
import "./App.css";

class App extends Component {
  constructor() {
    super();
    this.state = {
      psgState: {
        running: false
      }
    };
  }
  async componentDidMount() {
    let newState = {};
    // if (process.env.NODE_ENV === "PROD") {
    try {
      newState = await axios
        .get("/api/getState")
        .then(res => {
          console.log(res);
          if (res.status === 200) {
            return res.data;
          } else if (res.status === 429) {
            return { running: false, error: true, status: "reload" };
          } else if (res.status === 500) {
            return { running: false, error: true, status: "fatal" };
          } else {
            return { running: false, error: true, status: "fatal" };
          }
        })
        .catch(err => {
          console.error(err);
          return {
            running: false,
            error: true
          };
        });
    } catch (e) {
      console.error(e);
      newState = {
        running: false,
        error: true
      };
    }
    // } else {
    // newState = await new Promise((resolve, reject) => {
    //   resolve({
    //     songs: [
    //       { name: "First Song", artist: "First Artist" },
    //       { name: "Second Song", artist: "Second Artist" },
    //       { name: "Third Song", artist: "Third Artist" }
    //     ],
    //     running: true,
    //     currentSongIndex: 0,
    //     currentSong: { name: "First Song", artist: "First Artist" }
    //   });
    // });
    // }
    this.setState({ psgState: newState });
  }
  render() {
    let AppHead = () => (
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <h1 className="App-title">Welcome to Lightshow Pi</h1>
      </header>
    );
    if (this.state.psgState.error) {
      return (
        <div className="App">
          <AppHead />
          <p className="App-intro">An error occurred retrieving your state</p>
          {this.state.psgState.status &&
            this.state.psgState.status === "reload" && <p>Please Reload</p>}
        </div>
      );
    } else if (!this.state.psgState.running) {
      return (
        <div className="App">
          <AppHead />
          <p className="App-intro">The lightshow is currently unavailable</p>
        </div>
      );
    } else {
      return (
        <div className="App">
          <AppHead />
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
