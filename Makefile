CC = g++

all: main

main: restart.cpp
	$(CC) restart.cpp -o restart

clean:
	rm restart