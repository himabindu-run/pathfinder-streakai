import React, { useState } from "react";
import axios from "axios";
import "./Grid.css";

const Grid = () => {
  const [grid, setGrid] = useState(createGrid(5, 5));
  const [start, setStart] = useState(null);
  const [end, setEnd] = useState(null);
  const [path, setPath] = useState([]);

  const handleCellClick = (row, col) => {
    if (!start) {
      setStart([row, col]);
    } else if (!end) {
      setEnd([row, col]);
      findPath([row, col]);
    }
  };

  const findPath = async (selectedEnd) => {
    try {
      const response = await axios.post("http://localhost:8080/find-path", {
        start: { x: start[0], y: start[1] },
        end: { x: selectedEnd[0], y: selectedEnd[1] },
      });
      console.log("Path found:", response.data.path);
      setPath(response.data.path);
    } catch (error) {
      console.error("Error fetching path:", error);
    }
  };

  const renderCell = (row, col) => {
    const isStart = start && start[0] === row && start[1] === col;
    const isEnd = end && end[0] === row && end[1] === col;
    const isPath = path.some((p) => p.x === row && p.y === col);

    let className = "cell";
    if (isStart) className += " start";
    if (isEnd) className += " end";
    if (isPath) className += " path";

    return (
      <div
        key={`${row}-${col}`}
        className={className}
        onClick={() => handleCellClick(row, col)}
      />
    );
  };

  return (
    <div className="grid">
      {grid.map((row, rowIndex) => (
        <div key={rowIndex} className="row">
          {row.map((col, colIndex) => renderCell(rowIndex, colIndex))}
        </div>
      ))}
    </div>
  );
};

const createGrid = (rows, cols) => {
  return Array.from({ length: rows }, () => Array(cols).fill(0));
};

export default Grid;
