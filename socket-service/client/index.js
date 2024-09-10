const io = require("socket.io-client");
const args = process.argv.slice(2);

if (args.length === 0) {
  console.log("Usage: node client.js [0|1] (0 for consumer, 1 for producer)");
  process.exit(1);
}

const mode = args[0];
const roomID = "room-1234"; // The ID of the game room

// Connect to the Socket.IO server
const socket = io("http://localhost:3000", {
  transports: ["websocket", "polling"],
  reconnection: true,
  reconnectionAttempts: 5,
  timeout: 10000, 
  pingInterval: 25000, // Send a ping every 25 seconds
  pingTimeout: 60000, // Timeout if no pong response is received within 60 seconds
});

socket.on("connect", () => {
  console.log("Connected to the server");

  if (mode === "0") {
    // Consumer mode: Join the room and listen for events
    console.log("Joining as a consumer...");
    socket.emit("joinRoom", roomID);

    // Listen for events from the server
    socket.on("event", (data) => {
      console.log("Received event:", data);
    });

  } else if (mode === "1") {
    // Producer mode: Join the room and send events
    console.log("Joining as a producer...");
    socket.emit("joinRoom", roomID);

    // Send data every second
    setInterval(() => {
      const eventData = `Data from producer at ${new Date().toISOString()}`;
      socket.emit("produceEvent", roomID, eventData);
      console.log("Sent event:", eventData);
    }, 1000);
  } else {
    console.log("Invalid argument. Use 0 for consumer and 1 for producer.");
    process.exit(1);
  }
});

// Handle error events
socket.on("error", (err) => {
  console.error("Error:", err);
});

// Handle disconnection
socket.on("disconnect", () => {
  console.log("Disconnected from server.");
});
