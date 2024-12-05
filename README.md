
# Karaoke Backend

  

This project provides the backend services for the Karaoke Game where a user can join via their phone or TV (unity app) and play the game from a game room.

  

## Prerequisites

  

Ensure you have the following dependencies installed:

  

1.  **Docker Desktop**: Required to run containers for the backend and database services.

2.  **Makefile**: Used to streamline setup and deployment commands.

  

## Setup

  

Follow these steps to start the backend server.

  

### Step 1: Start the Backend

  

To build and run the backend, execute the following command from the `/start` folder:

  

```bash

make  up_build

```

  

To shut it down:

  

```bash

make  down

```

## How to operate Local Mode
Make a POST request at:
`http://localhost:8080/games/local`
With the following body:
`{
"minPlayers":  2,
"maxPlayers":  10,
"createdBy":  "<username/id>"
}`

This will create a stream within Redis and set up our database. The game room gets created and users can join, and play within the gameroom until the events stop.

## TODO
- Matchmaking
- Live Game mode integration with matchmaking

## Sample Scripts
To run sample client scripts, go to `game-service/client`

To create a producer (represented by: 1 and sends events i.e in our case is the mobile user):
`go run main.go 1`

To create a receiver (represented by: 0 and can only receive events i.e, in our case is the television user):
`go run main.go 0`