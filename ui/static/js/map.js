const mapElem = document.getElementById("map");
const mapWidth = mapElem.clientWidth;
const mapHeight = mapElem.clientHeight;
const moveRobot = (robot, movementScaleFactor) => {
  let robotElement = document.getElementById(robot.Name);
  if (!robotElement) {
    robotElement = document.createElement("div");
    robotElement.id = robot.Name;
    robotElement.className = "robot";
    mapElem.appendChild(robotElement);
  }

  const robotSize = 20;
  robotElement.style.width = `${robotSize}px`;
  robotElement.style.height = `${robotSize}px`;

  const robotCenterX = robot.Location.X + robotSize / 2;
  const robotCenterY = robot.Location.Y + robotSize / 2;

  const robotX = movementScaleFactor * robotCenterX + mapWidth / 2;
  const robotY = movementScaleFactor * robotCenterY + mapHeight / 2;
  robotElement.style.left = `${robotX}px`;
  robotElement.style.top = `${robotY}px`;
};

(() => {
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

  const topic = "robots/locations";
  const client = mqtt.connect("ws://localhost:8083");
  client.on("connect", () => {
    client.subscribe(topic);
  });
  client.on("message", (_topic, message) => {
    if (_topic !== topic) return;
    const data = JSON.parse(message.toString());
    console.log(data);
    const robot = robots.find((robot) => robot?.Name === data?.Name);
    if (robot) {
      robot.Location = data.Location;
    } else {
      robots.push(data);
    }
    moveRobot(data, 10);
  });
})();
