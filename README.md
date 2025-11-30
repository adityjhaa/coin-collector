# Coin Collector 
**Multiplayer Game (Go + Ebiten + UDP)**
_Authoritative Server • Interpolation • Lag Simulation • Smooth Gameplay_

---

## 🎮 Overview

**Coin Collector** is a multiplayer 2D game built using **Go**, **Ebiten**, and **UDP networking**.  
Players join a server, move around a shared arena, collect coins, and score points — all under strict **server authority**, simulated **network latency**, and smooth **interpolation-based rendering**.

---

## 🏗️ Project Structure

```
coin-collector/
│
├── common/         # Shared protocol, constants, types
├── server/         # Game server logic (state, coins, players)
├── client/         # Rendering, interpolation, networking
└── cmd/
    ├── server/     # server main entry
    └── client/     # client main entry
```

---

## 📚 Dependencies
- Go 1.21+
- Ebiten v2

---

## 🚀 Running the Game

### 0️⃣ Clone the repository

```bash
git clone https://github.com/adityjhaa/coin-collector.git
cd coin-collector
```


### 1️⃣ Install dependencies

```bash
go mod tidy
```

### 2️⃣ Run the server

```bash
go run cmd/server/main.go
```

### 3️⃣ Run two clients (each in a separate terminal or window)

```bash
go run cmd/client/main.go
```

---

## 🧠 Technical Highlights

### Interpolation
Clients render slightly behind server time (`InterpDelayMs = 200ms`) so they always have two snapshots to interpolate between.  
This prevents jitter and keeps remote players smooth even under network latency.

### Input-Only Client
Clients send only a movement bitmask (W/A/S/D).  
No position, velocity, or physics data ever comes from the client.

### Server Authority
All simulation is performed on the server:

- Movement  
- Collisions  
- Coin spawning  
- Scoring  
- Player lifecycle  

Clients only predict **input**, not state.

### Latency Simulation
All outgoing packets are delayed by `SimulatedLatencyMs` to emulate real online conditions  
and to demonstrate robust interpolation handling.

---

## 🧑‍💻 Author
**Aditya Jha**
