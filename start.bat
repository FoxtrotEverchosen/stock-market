@echo off
set PORT=%1
if "%PORT%"=="" set PORT=8080
docker compose up --build