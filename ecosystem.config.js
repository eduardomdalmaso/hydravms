module.exports = {
  apps: [
    // 1. HydraVMS Backend - Control Plane REST API & WebSockets (porta 8083)
    {
      name: "hydra-vms-api",
      cwd: "/home/hades/Documents/HydraVMS",
      script: "./bin/hydravms",
      interpreter: "none",
      autorestart: true,
      max_restarts: 10,
      restart_delay: 2000,
      env: {
        PORT: "8083",
        DATABASE_URL: "postgres://postgres:postgres@127.0.0.1:5432/hydravms?sslmode=disable",
        NATS_URL: "nats://127.0.0.1:4222",
        MINIO_ENDPOINT: "127.0.0.1:9000",
      },
    },

    // 2. HydraVMS Frontend - Painel de Controle e Orquestrador Web (porta 5173)
    {
      name: "hydra-vms",
      cwd: "/home/hades/Documents/HydraVMS/web",
      script: "npm",
      args: "run dev",
      autorestart: true,
      max_restarts: 10,
      restart_delay: 2000,
      env: {
        PORT: "5173",
      },
    },
  ],
};
