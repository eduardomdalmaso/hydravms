const path = require("path");
const fs = require("fs");

const apps = [
  // 1. HydraVMS Backend - Control Plane REST API & WebSockets (porta 8083)
  {
    name: "hydra-vms-api",
    cwd: __dirname,
    script: process.platform === "win32" ? "./hydravms.exe" : "./bin/hydravms",
    interpreter: "none",
    autorestart: true,
    max_restarts: 10,
    restart_delay: 2000,
    env: {
      PORT: "8083",
      DB_DRIVER: "sqlite",
      SQLITE_PATH: "hydravms.db",
      JWT_SECRET: "hydravms-enterprise-secure-jwt-signing-key-32b-secret",
      NATS_URL: "nats://127.0.0.1:4222",
      MINIO_ENDPOINT: "127.0.0.1:9000",
    },
  },

  // 2. HydraVMS Frontend - Painel de Controle e Orquestrador Web (porta 5173)
  {
    name: "hydra-vms",
    cwd: path.join(__dirname, "web"),
    script: "./node_modules/vite/bin/vite.js",
    args: "--host --port 5173",
    autorestart: true,
    max_restarts: 10,
    restart_delay: 2000,
    env: {
      PORT: "5173",
    },
  },

  // 3. HydraStream - Motor de Mídia, Probing e Descoberta ONVIF WS-Discovery (porta 8080)
  {
    name: "hydra-stream",
    cwd: path.resolve(__dirname, "..", "HydraStream"),
    script: process.platform === "win32" ? "./hydrastream.exe" : "./bin/hydrastream",
    interpreter: "none",
    autorestart: true,
    max_restarts: 10,
    restart_delay: 2000,
    env: {
      PORT: "8080",
    },
  },
];

// 4. HydraForge (Opcional - se binário existir)
const forgeDir = path.resolve(__dirname, "..", "HydraForge");
const forgeBin = process.platform === "win32" ? path.join(forgeDir, "hydraforge.exe") : path.join(forgeDir, "bin", "hydraforge");
if (fs.existsSync(forgeBin)) {
  apps.push({
    name: "hydra-forge",
    cwd: forgeDir,
    script: forgeBin,
    interpreter: "none",
    autorestart: true,
    max_restarts: 10,
    restart_delay: 2000,
    env: {
      PORT: "8081",
      JWT_SECRET: "hydravms-enterprise-secure-jwt-signing-key-32b-secret",
    },
  });
}

// 5. HydraVault (Opcional - se binário existir)
const vaultDir = path.resolve(__dirname, "..", "HydraVault");
const vaultBin = process.platform === "win32" ? path.join(vaultDir, "hydravault.exe") : path.join(vaultDir, "bin", "hydravault");
if (fs.existsSync(vaultBin)) {
  apps.push({
    name: "hydra-vault",
    cwd: vaultDir,
    script: vaultBin,
    interpreter: "none",
    autorestart: true,
    max_restarts: 10,
    restart_delay: 2000,
    env: {
      PORT: "8082",
    },
  });
}

module.exports = { apps };

