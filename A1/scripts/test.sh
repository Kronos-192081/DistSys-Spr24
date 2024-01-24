#!/bin/bash
echo "Launching 10000 async requests on http://localhost:5000/home route via python3"
cd ./Analysis
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python3 analysis.py
deactivate
rm -rf venv
echo "Statistics generated."
xdg-open A1.png
xdg-open A2.png