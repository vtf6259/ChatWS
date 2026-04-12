#!/bin/sh
cd frontend
npm i
npm run build
mkdir ../build/dist
cp -R dist ../build/dist/ChatWS