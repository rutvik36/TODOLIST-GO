# pkg

This directory is reserved for **public, reusable packages** that could be imported by other modules (for example, shared API types or clients).

The to-do server keeps its implementation under `internal/` so the application boundary stays clear. If you later extract a client SDK or shared validators, place them here.
