import React, { useEffect, useState } from "react";
import './App.css';

function renderSVG(layout) {
  if (!layout) return null;
  const { glines = [], stations = [], stationareas = [], trainzones = [], blockzones = [], buttons = [], labels = [] } = layout;
  return (
    <svg id="map" className="border" width="100%" viewBox="0 0 1000 1300" xmlns="http://www.w3.org/2000/svg">
      <g id="layer">
        {glines.map((line, i) => (
          <line key={i} x1={line.x1} y1={line.y1} x2={line.x2} y2={line.y2} stroke={line.stroke} strokeWidth={line.stroke_width} />
        ))}
        {stations.map((stn, i) => (
          <circle key={i} cx={stn.x} cy={stn.y} r={12} fill="dodgerblue" />
        ))}
        {stations.map((stn, i) => (
          <text key={"lbl"+i} x={stn.x+16} y={stn.y} fontSize="18" fill="#fff">{stn.name}</text>
        ))}
        {/* Add more SVG elements for stationareas, trainzones, blockzones, buttons, labels as needed */}
      </g>
    </svg>
  );
}

function App() {
  const [servers, setServers] = useState([]);
  const [layouts, setLayouts] = useState([]);
  const [server, setServer] = useState("en1");
  const [layoutNr, setLayoutNr] = useState(0);
  const [layout, setLayout] = useState(null);
  const [trains, setTrains] = useState([]);
  const [selectedTrain, setSelectedTrain] = useState(null);

  useEffect(() => {
    fetch(`/api?type=layout&server=${server}&layout=${layoutNr}`)
      .then((res) => res.json())
      .then((data) => setLayout(data.layout_data));
  }, [server, layoutNr]);

  useEffect(() => {
    fetch(`/api?type=trains&server=${server}&layout=${layoutNr}`)
      .then((res) => res.json())
      .then((data) => setTrains(data.data?.t || []));
  }, [server, layoutNr]);

  useEffect(() => {
    fetch('/api?type=servers')
      .then(res => res.json())
      .then(data => setServers(data.servers || []));
    fetch('/api?type=layouts')
      .then(res => res.json())
      .then(data => setLayouts(data.layouts || []));
  }, []);

  return (
    <div className="container">
      <header className="d-flex py-3">
        <ul className="nav nav-pills">
          <li className="nav-item"><span className="logo">Simrail Dispatcher Eye (SiDE)</span></li>
          <li className="nav-item"><a href="#" className="nav-link">Home</a></li>
          <li className="nav-item"><a href="#" className="nav-link">Download</a></li>
        </ul>
      </header>
      <div className="row gx-5 py-3 mx-0">
        <div className="col card border-0 m-1">
          <h2>Map (Layout {layoutNr})</h2>
          <div style={{ minHeight: 300, background: '#222', color: '#fff', padding: 10 }}>
            {renderSVG(layout)}
          </div>
          <div>
            <label>Change layout: </label>
            <select value={layoutNr} onChange={e => setLayoutNr(Number(e.target.value))}>
              {layouts.map(l => <option key={l.id} value={l.number}>{l.name}</option>)}
            </select>
          </div>
        </div>
        <div className="col card border">
          <h2>Train info</h2>
          <div>
            <label>Change server: </label>
            <select value={server} onChange={e => setServer(e.target.value)}>
              {servers.map(s => <option key={s.ServerCode} value={s.ServerCode}>{s.ServerName}</option>)}
            </select>
          </div>
          <table className="table table-striped table-hover">
            <thead>
              <tr>
                <th>Train No</th>
                <th>Name</th>
                <th>Speed</th>
                <th>Delay</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              {trains.map(train => (
                <tr key={train.trainNoLocal}>
                  <td>{train.trainNoLocal}</td>
                  <td>{train.trainObject?.TrainName}</td>
                  <td>{train.trainObject?.TrainData?.Velocity}</td>
                  <td>{train.trainObject?.delay}</td>
                  <td><button onClick={() => setSelectedTrain(train)}>Details</button></td>
                </tr>
              ))}
            </tbody>
          </table>
          {selectedTrain && (
            <div className="container">
              <h3>Train Details</h3>
              <div><strong>No:</strong> {selectedTrain.trainNoLocal}</div>
              <div><strong>Name:</strong> {selectedTrain.trainObject?.TrainName}</div>
              <div><strong>Speed:</strong> {selectedTrain.trainObject?.TrainData?.Velocity}</div>
              <div><strong>Delay:</strong> {selectedTrain.trainObject?.delay}</div>
              <div><strong>Locotype:</strong> {selectedTrain.trainObject?.Locotype}</div>
              <div><strong>Length:</strong> {selectedTrain.trainObject?.TrainLength}</div>
              <div><strong>Weight:</strong> {selectedTrain.trainObject?.TrainWeight}</div>
              <div><strong>Type:</strong> {selectedTrain.trainObject?.TrainType}</div>
              <div><strong>Start Station:</strong> {selectedTrain.trainObject?.startStation}</div>
              <div><strong>End Station:</strong> {selectedTrain.trainObject?.endStation}</div>
              <div><strong>Held:</strong> {selectedTrain.trainObject?.held_at_signal ? "Yes" : "No"}</div>
              <div><strong>Held Soon:</strong> {selectedTrain.trainObject?.held_at_signal_soon ? "Yes" : "No"}</div>
              <h4>Timetable</h4>
              <table className="table table-striped table-hover">
                <thead>
                  <tr>
                    <th>Station</th>
                    <th>Arr</th>
                    <th>Dep</th>
                    <th>Stop</th>
                  </tr>
                </thead>
                <tbody>
                  {(selectedTrain.timetable || []).map((tt, idx) => (
                    <tr key={idx}>
                      <td>{tt.nameForPerson}</td>
                      <td>{tt.arrivalTime_real || tt.arrivalTime}</td>
                      <td>{tt.departureTime_real || tt.departureTime}</td>
                      <td>{tt.stop_duration}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default App;
