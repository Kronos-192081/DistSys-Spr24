#!/bin/bash
# echo "Launching 10000 async requests on http://localhost:5000/ route via python3"
echo -ne "Enter 1 for asynchronous testing or 2 for synchronous: "
read test_type
echo -e "Setting up Python virtual environment"

cd ./Analysis
python3 -m venv .venv
source .venv/bin/activate

pip install requests aiohttp asyncio >/dev/null

if [ "$test_type" = "1" ]; then
	echo -e "\nRunning test 1"
	python3 analysis_async.py 1

	echo -e "\nRunning test 2"
	python3 analysis_async.py 2

	echo -e "\nRunning test 3"
	python3 analysis_async.py 3

elif [ "$test_type" = "2" ]; then
	echo -e "\nRunning test 1"
	python3 analysis_sync.py 1

	echo -e "\nRunning test 2"
	python3 analysis_sync.py 2

	echo -e "\nRunning test 3"
	python3 analysis_sync.py 3
else
    echo "Invalid input. Please enter 1 or 2."
fi

deactivate
rm -rf .venv

echo -e "Tests run successfully. Execute 'make kill' and 'make clean' in that order."