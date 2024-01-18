/**
 * @type {Robot[]}
 * @typedef Robot
 * @property {string} Name
 * @property {Location} Location
 *
 * @typedef Location
 * @property {number} X
 * @property {number} Y
 */
const robots = [];

/**
 * @type {Map}
 * @typedef Map
 * @property {Line[]} Lines
 * @property {Size} Size
 *
 * @typedef Line
 * @property {Point} Start
 * @property {Point} End
 *
 * @typedef Point
 * @property {number} X
 * @property {number} Y
 *
 * @typedef Size
 * @property {number} Width
 * @property {number} Height
 */
let map = {};

const mapCanvas = document.getElementById("map");

(() => {
  const robotsLocationTopic = "robots/locations";
  const mapsFloorTopic = "maps/floor";

  const client = mqtt.connect("ws://localhost:8083");
  client.on("connect", () => {
    client.subscribe(robotsLocationTopic);
    client.subscribe(mapsFloorTopic);
  });
  client.on("message", (_topic, message) => {
    switch (_topic) {
      case robotsLocationTopic:
        // const data = JSON.parse(message.toString());
        // const robot = robots.find((robot) => robot?.Name === data?.Name);
        // if (robot) {
        //   robot.Location = data.Location;
        // } else {
        //   robots.push(data);
        // }
        // moveRobot(data);
        break;
      case mapsFloorTopic:
        map = JSON.parse(message.toString());
        drawMap(map);
        // mapElem.style.backgroundImage = `url(${map.Image})`;
        break;
      default:
        break;
    }
  });
})();

const moveRobot = (robot) => {
  // let robotElement = document.getElementById(robot.Name);
  // if (!robotElement) {
  //   robotElement = document.createElement("div");
  //   robotElement.id = robot.Name;
  //   robotElement.className = "robot";
  //   mapElem.appendChild(robotElement);
  // }
  // const robotSize = 20;
  // robotElement.style.width = `${robotSize}px`;
  // robotElement.style.height = `${robotSize}px`;
  // robotElement.style.position = "absolute";
  // robotElement.style.left = `${robot.Location.X}px`;
  // robotElement.style.top = `${robot.Location.Y}px`;
  // robotElement.style.backgroundColor = "white";
};

const drawMap = (map) => {
  mapCanvas.width = map.Size.Width;
  mapCanvas.height = map.Size.Height;

  const ctx = mapCanvas.getContext("2d");

  ctx.clearRect(0, 0, map.Size.Width, map.Size.Height);

  ctx.fillStyle = "black";
  ctx.fillRect(0, 0, map.Size.Width, map.Size.Height);

  ctx.beginPath();
  ctx.lineWidth = 5;
  ctx.strokeStyle = "red";

  map.Lines.forEach((line) => {
    ctx.moveTo(line.Start.X, line.Start.Y);
    ctx.lineTo(line.End.X, line.End.Y);
  });

  ctx.stroke();
};
