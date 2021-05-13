SHELL := /bin/bash

#=============================================================

#Speedy is a Go apllication API that calculates and show your internet speed from 
#fast.com and speedtest.com APIs together in a single 
#application without having to open two tabs to verify 
#your internet speed, When a client goto:http://localhost:3000/v1/getspeed
#it requests for the internet speed and respond back with a JSON response of
#the internet speed.

#NB: Test cases has not been written yet for this project.

#==============================================================	

tidy:
	go mod tidy
	go mod vendor