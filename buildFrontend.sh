#!/bin/sh
cd frontend
npm i
npm run build
cp -R dist ../build